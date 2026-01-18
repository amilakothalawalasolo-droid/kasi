package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var db *sql.DB
var store *sessions.CookieStore

type Expense struct { ID int; Item string; Amount float64; Type string; Category string; Quantity float64; Unit string; Username string; Date string; UserID int }
type User struct { ID int; Username string; IsAdmin bool; Currency string; Language string; ProjectName string }
type Category struct { ID int; Name string }
type Unit struct { ID int; Name string }
type ReportItem struct { Category string; Total float64 }

type PageData struct {
	User User; Private []Expense; Common []Expense; Categories []Category; Units []Unit;
	PrivTot float64; ComTot float64; MonthPriv float64; MonthCom float64;
	CSRFField template.HTML; Error string; Success string; FilterStart string; FilterEnd string; Today string; UsersList []User;
	ReportSummary []ReportItem; ReportTotal float64;
}

func main() {
	if _, err := os.Stat("./data"); os.IsNotExist(err) { os.Mkdir("./data", 0755) }
	secretKey := os.Getenv("SESSION_SECRET"); if secretKey == "" { secretKey = "dev-secret" }
	store = sessions.NewCookieStore([]byte(secretKey)); store.Options.HttpOnly = true; store.Options.SameSite = http.SameSiteLaxMode; store.Options.Secure = false
	initDBConnection(); initDB()
	csrfMiddleware := csrf.Protect([]byte(secretKey), csrf.Secure(false))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/", authMiddleware(homeHandler)); mux.HandleFunc("/login", loginHandler); mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/add", authMiddleware(addHandler)); mux.HandleFunc("/delete", authMiddleware(deleteHandler)); mux.HandleFunc("/edit", authMiddleware(editHandler))
	mux.HandleFunc("/settings", authMiddleware(settingsHandler))
	mux.HandleFunc("/settings/category/add", authMiddleware(addCategoryHandler)); mux.HandleFunc("/settings/category/delete", authMiddleware(deleteCategoryHandler))
	mux.HandleFunc("/settings/unit/add", authMiddleware(addUnitHandler)); mux.HandleFunc("/settings/unit/delete", authMiddleware(deleteUnitHandler))
	mux.HandleFunc("/report", authMiddleware(reportHandler))
	mux.HandleFunc("/admin", authMiddleware(adminHandler)); mux.HandleFunc("/admin/create", authMiddleware(adminCreateUserHandler)); mux.HandleFunc("/admin/delete", authMiddleware(adminDeleteUserHandler)); mux.HandleFunc("/admin/edit", authMiddleware(adminEditUserHandler))
	mux.HandleFunc("/admin/backup", authMiddleware(adminBackupHandler)); mux.HandleFunc("/admin/restore", authMiddleware(adminRestoreHandler))

	log.Println("🪙 Kasi v5.0.1 (Stable) started on :8080")
	http.ListenAndServe(":8080", csrfMiddleware(mux))
}

func initDBConnection() { var err error; if db != nil { db.Close() }; db, err = sql.Open("sqlite", "./data/kasi.db"); if err != nil { log.Fatal(err) } }

func initDB() {
	db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE, password TEXT, is_admin BOOLEAN DEFAULT 0, currency TEXT DEFAULT 'LKR', language TEXT DEFAULT 'en', project_name TEXT DEFAULT 'Project Expenses');`)
	db.Exec(`CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE);`)
	db.Exec(`CREATE TABLE IF NOT EXISTS units (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE);`)
	db.Exec(`CREATE TABLE IF NOT EXISTS expenses (id INTEGER PRIMARY KEY AUTOINCREMENT, item TEXT, amount REAL, type TEXT, category TEXT, quantity REAL DEFAULT 1, unit TEXT, user_id INTEGER, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY(user_id) REFERENCES users(id));`)
	
	// --- FIX START: Check total users count instead of just 'admin' ---
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 { 
		// Only create default admin if NO users exist at all
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
		db.Exec("INSERT INTO users(username, password, is_admin, currency, language, project_name) VALUES(?, ?, ?, ?, ?, ?)", "admin", string(hash), true, "LKR", "en", "My Project") 
	}
	// --- FIX END ---

	db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count); if count == 0 { db.Exec("INSERT INTO categories (name) VALUES ('Food'), ('Transport'), ('Materials'), ('Labor'), ('Other')") }
	db.QueryRow("SELECT COUNT(*) FROM units").Scan(&count); if count == 0 { db.Exec("INSERT INTO units (name) VALUES ('Items'), ('kg'), ('L'), ('Day'), ('Feet')") }
}

func getUser(id int) User { var u User; db.QueryRow("SELECT id, username, is_admin, currency, language, COALESCE(project_name, 'Project') FROM users WHERE id=?", id).Scan(&u.ID, &u.Username, &u.IsAdmin, &u.Currency, &u.Language, &u.ProjectName); return u }

func authMiddleware(next http.HandlerFunc) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) { session, _ := store.Get(r, "kasi-session"); if val, ok := session.Values["userID"]; !ok || val == nil { http.Redirect(w, r, "/login", http.StatusSeeOther); return }; next(w, r) } }
func loginHandler(w http.ResponseWriter, r *http.Request) { if r.Method == "GET" { template.Must(template.ParseFiles("templates/login.html")).Execute(w, map[string]interface{}{"CSRFField": csrf.TemplateField(r), "Error": r.URL.Query().Get("error")}); return }; username := r.FormValue("username"); password := r.FormValue("password"); var id int; var hash string; err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&id, &hash); if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil { http.Redirect(w, r, "/login?error=Invalid Credentials", http.StatusSeeOther); return }; session, _ := store.Get(r, "kasi-session"); session.Values["userID"] = id; session.Save(r, w); http.Redirect(w, r, "/", http.StatusSeeOther) }
func logoutHandler(w http.ResponseWriter, r *http.Request) { session, _ := store.Get(r, "kasi-session"); session.Values["userID"] = nil; session.Options.MaxAge = -1; session.Save(r, w); http.Redirect(w, r, "/login", http.StatusSeeOther) }

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); currentUser := getUser(userID)
	var monthPriv, monthCom float64; currentMonth := time.Now().Format("2006-01")
	db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE user_id=? AND type='private' AND strftime('%Y-%m', created_at) = ?", userID, currentMonth).Scan(&monthPriv)
	db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE type='common' AND strftime('%Y-%m', created_at) = ?", currentMonth).Scan(&monthCom)
	start := r.URL.Query().Get("start"); end := r.URL.Query().Get("end"); filterQuery := ""; args := []interface{}{}; if start != "" && end != "" { filterQuery = " AND date(created_at) BETWEEN ? AND ?"; args = append(args, start, end) }
	catRows, _ := db.Query("SELECT id, name FROM categories ORDER BY name ASC"); var categories []Category; for catRows.Next() { var c Category; catRows.Scan(&c.ID, &c.Name); categories = append(categories, c) }; catRows.Close()
	unitRows, _ := db.Query("SELECT id, name FROM units ORDER BY name ASC"); var units []Unit; for unitRows.Next() { var u Unit; unitRows.Scan(&u.ID, &u.Name); units = append(units, u) }; unitRows.Close()
	privQuery := "SELECT id, item, amount, type, category, quantity, unit, strftime('%Y-%m-%d', created_at), user_id FROM expenses WHERE user_id = ? AND type = 'private'" + filterQuery + " ORDER BY created_at DESC LIMIT 50"; privArgs := append([]interface{}{userID}, args...); rowsPriv, _ := db.Query(privQuery, privArgs...); var privateExpenses []Expense; var privTotal float64; if rowsPriv != nil { for rowsPriv.Next() { var e Expense; rowsPriv.Scan(&e.ID, &e.Item, &e.Amount, &e.Type, &e.Category, &e.Quantity, &e.Unit, &e.Date, &e.UserID); privateExpenses = append(privateExpenses, e); privTotal += e.Amount }; rowsPriv.Close() }
	comQuery := "SELECT e.id, e.item, e.amount, e.type, e.category, e.quantity, e.unit, strftime('%Y-%m-%d', e.created_at), u.username, e.user_id FROM expenses e JOIN users u ON e.user_id = u.id WHERE e.type = 'common'" + filterQuery + " ORDER BY e.created_at DESC LIMIT 50"; rowsCom, _ := db.Query(comQuery, args...); var commonExpenses []Expense; var comTotal float64; if rowsCom != nil { for rowsCom.Next() { var e Expense; rowsCom.Scan(&e.ID, &e.Item, &e.Amount, &e.Type, &e.Category, &e.Quantity, &e.Unit, &e.Date, &e.Username, &e.UserID); commonExpenses = append(commonExpenses, e); comTotal += e.Amount }; rowsCom.Close() }
	tmpl := template.Must(template.ParseFiles("templates/home.html")); data := PageData{User: currentUser, Private: privateExpenses, Common: commonExpenses, Categories: categories, Units: units, PrivTot: privTotal, ComTot: comTotal, MonthPriv: monthPriv, MonthCom: monthCom, CSRFField: csrf.TemplateField(r), FilterStart: start, FilterEnd: end, Today: time.Now().Format("2006-01-02"), Error: r.URL.Query().Get("error"), Success: r.URL.Query().Get("success")}; tmpl.Execute(w, data)
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); currentUser := getUser(userID)
	start := r.URL.Query().Get("start"); end := r.URL.Query().Get("end"); filterQuery := ""; args := []interface{}{}; if start != "" && end != "" { filterQuery = " AND date(created_at) BETWEEN ? AND ?"; args = append(args, start, end) }
	catQuery := "SELECT category, COALESCE(SUM(amount), 0) FROM expenses WHERE type='common'" + filterQuery + " GROUP BY category ORDER BY SUM(amount) DESC"; rows, err := db.Query(catQuery, args...); var summary []ReportItem; var grandTotal float64
	if err == nil { for rows.Next() { var ri ReportItem; rows.Scan(&ri.Category, &ri.Total); summary = append(summary, ri); grandTotal += ri.Total }; rows.Close() }
	listQuery := "SELECT e.item, e.amount, e.category, e.date, u.username FROM expenses e JOIN users u ON e.user_id = u.id WHERE e.type='common'" + filterQuery + " ORDER BY e.created_at DESC"; rowsList, errList := db.Query(listQuery, args...); var details []Expense
	if errList == nil { for rowsList.Next() { var e Expense; rowsList.Scan(&e.Item, &e.Amount, &e.Category, &e.Date, &e.Username); details = append(details, e) }; rowsList.Close() }
	tmpl, errTmpl := template.ParseFiles("templates/report.html")
	if errTmpl != nil { http.Error(w, "Report template missing", 500); return }
	data := PageData{User: currentUser, ReportSummary: summary, ReportTotal: grandTotal, Common: details, FilterStart: start, FilterEnd: end, Today: time.Now().Format("2006-01-02")}; tmpl.Execute(w, data)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int)
	if r.Method == "POST" {
		username := r.FormValue("username"); currency := r.FormValue("currency"); language := r.FormValue("language"); password := r.FormValue("password"); projectName := r.FormValue("project_name")
		_, err := db.Exec("UPDATE users SET username=?, currency=?, language=? WHERE id=?", username, currency, language, userID); if err != nil { http.Redirect(w, r, "/settings?error=Username taken", http.StatusSeeOther); return }
		if projectName != "" { db.Exec("UPDATE users SET project_name=?", projectName) }
		if password != "" { hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10); db.Exec("UPDATE users SET password=? WHERE id=?", string(hashed), userID) }
		http.Redirect(w, r, "/settings?success=Saved & Synced", http.StatusSeeOther); return
	}
	catRows, _ := db.Query("SELECT id, name FROM categories ORDER BY name ASC"); var categories []Category; for catRows.Next() { var c Category; catRows.Scan(&c.ID, &c.Name); categories = append(categories, c) }; catRows.Close(); unitRows, _ := db.Query("SELECT id, name FROM units ORDER BY name ASC"); var units []Unit; for unitRows.Next() { var u Unit; unitRows.Scan(&u.ID, &u.Name); units = append(units, u) }; unitRows.Close(); tmpl := template.Must(template.ParseFiles("templates/settings.html")); data := PageData{User: getUser(userID), Categories: categories, Units: units, CSRFField: csrf.TemplateField(r), Success: r.URL.Query().Get("success"), Error: r.URL.Query().Get("error")}; tmpl.Execute(w, data)
}

func addHandler(w http.ResponseWriter, r *http.Request) { if r.Method == "POST" { session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); item := r.FormValue("item"); amount := r.FormValue("amount"); expType := r.FormValue("type"); category := r.FormValue("category"); quantity := r.FormValue("quantity"); unit := r.FormValue("unit"); date := r.FormValue("date"); if date == "" { date = time.Now().Format("2006-01-02") }; if quantity == "" { quantity = "1" }; db.Exec("INSERT INTO expenses(item, amount, type, category, quantity, unit, user_id, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?)", item, amount, expType, category, quantity, unit, userID, date) }; http.Redirect(w, r, "/", http.StatusSeeOther) }
func addCategoryHandler(w http.ResponseWriter, r *http.Request) { if r.Method == "POST" { name := r.FormValue("name"); if name != "" { _, err := db.Exec("INSERT INTO categories(name) VALUES(?)", name); if err != nil { http.Redirect(w, r, "/settings?error=Exists", http.StatusSeeOther); return } } }; http.Redirect(w, r, "/settings?success=Added", http.StatusSeeOther) }
func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) { if r.Method == "POST" { id := r.FormValue("id"); db.Exec("DELETE FROM categories WHERE id=?", id) }; http.Redirect(w, r, "/settings?success=Deleted", http.StatusSeeOther) }
func addUnitHandler(w http.ResponseWriter, r *http.Request) { if r.Method == "POST" { name := r.FormValue("name"); if name != "" { _, err := db.Exec("INSERT INTO units(name) VALUES(?)", name); if err != nil { http.Redirect(w, r, "/settings?error=Exists", http.StatusSeeOther); return } } }; http.Redirect(w, r, "/settings?success=Added", http.StatusSeeOther) }
func deleteUnitHandler(w http.ResponseWriter, r *http.Request) { if r.Method == "POST" { id := r.FormValue("id"); db.Exec("DELETE FROM units WHERE id=?", id) }; http.Redirect(w, r, "/settings?success=Deleted", http.StatusSeeOther) }
func deleteHandler(w http.ResponseWriter, r *http.Request) { if r.Method == "POST" { id := r.FormValue("id"); session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); var isAdmin bool; db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); if isAdmin { db.Exec("DELETE FROM expenses WHERE id = ?", id) } else { db.Exec("DELETE FROM expenses WHERE id = ? AND user_id = ?", id, userID) } }; http.Redirect(w, r, "/", http.StatusSeeOther) }
func editHandler(w http.ResponseWriter, r *http.Request) { if r.Method == "POST" { id := r.FormValue("id"); item := r.FormValue("item"); amount := r.FormValue("amount"); session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); var isAdmin bool; db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); if isAdmin { db.Exec("UPDATE expenses SET item=?, amount=? WHERE id=?", item, amount, id) } else { db.Exec("UPDATE expenses SET item=?, amount=? WHERE id=? AND user_id=?", item, amount, id, userID) } }; http.Redirect(w, r, "/", http.StatusSeeOther) }
func adminHandler(w http.ResponseWriter, r *http.Request) { session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); var isAdmin bool; db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); if !isAdmin { http.Redirect(w, r, "/", http.StatusSeeOther); return }; rows, _ := db.Query("SELECT id, username, is_admin, currency, language FROM users"); var users []User; for rows.Next() { var u User; rows.Scan(&u.ID, &u.Username, &u.IsAdmin, &u.Currency, &u.Language); users = append(users, u) }; tmpl := template.Must(template.ParseFiles("templates/admin.html")); data := PageData{User: getUser(userID), UsersList: users, CSRFField: csrf.TemplateField(r), Error: r.URL.Query().Get("error"), Success: r.URL.Query().Get("success")}; tmpl.Execute(w, data) }
func adminCreateUserHandler(w http.ResponseWriter, r *http.Request) { session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); var isAdmin bool; db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); if !isAdmin { http.Redirect(w, r, "/", http.StatusSeeOther); return }; username := r.FormValue("username"); password := r.FormValue("password"); hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10); _, err := db.Exec("INSERT INTO users(username, password, currency, language, project_name) VALUES(?, ?, 'LKR', 'en', 'My Project')", username, string(hashed)); if err != nil { http.Redirect(w, r, "/admin?error=Username exists", http.StatusSeeOther); return }; http.Redirect(w, r, "/admin?success=User Created", http.StatusSeeOther) }
func adminDeleteUserHandler(w http.ResponseWriter, r *http.Request) { session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); var isAdmin bool; db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); if !isAdmin { http.Redirect(w, r, "/", http.StatusSeeOther); return }; id := r.FormValue("id"); if id == "1" { http.Redirect(w, r, "/admin?error=Cannot delete main admin", http.StatusSeeOther); return }; db.Exec("DELETE FROM expenses WHERE user_id = ?", id); db.Exec("DELETE FROM users WHERE id = ?", id); http.Redirect(w, r, "/admin?success=User Deleted", http.StatusSeeOther) }
func adminEditUserHandler(w http.ResponseWriter, r *http.Request) { session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); var isAdmin bool; db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); if !isAdmin { http.Redirect(w, r, "/", http.StatusSeeOther); return }; if r.Method == "POST" { targetID := r.FormValue("id"); currency := r.FormValue("currency"); password := r.FormValue("password"); db.Exec("UPDATE users SET currency=? WHERE id=?", currency, targetID); if password != "" { hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10); db.Exec("UPDATE users SET password=? WHERE id=?", string(hashed), targetID) }; http.Redirect(w, r, "/admin?success=User Updated", http.StatusSeeOther) } }
func adminBackupHandler(w http.ResponseWriter, r *http.Request) { session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); var isAdmin bool; db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); if !isAdmin { http.Error(w, "Unauthorized", http.StatusUnauthorized); return }; w.Header().Set("Content-Disposition", "attachment; filename=kasi_backup_"+time.Now().Format("20060102")+".db"); w.Header().Set("Content-Type", "application/octet-stream"); http.ServeFile(w, r, "./data/kasi.db") }
func adminRestoreHandler(w http.ResponseWriter, r *http.Request) { session, _ := store.Get(r, "kasi-session"); userID := session.Values["userID"].(int); var isAdmin bool; db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); if !isAdmin { http.Redirect(w, r, "/", http.StatusSeeOther); return }; if r.Method == "POST" { file, _, err := r.FormFile("backup_file"); if err != nil { http.Redirect(w, r, "/admin?error=File Error", http.StatusSeeOther); return }; defer file.Close(); db.Close(); out, err := os.Create("./data/kasi.db"); if err != nil { initDBConnection(); http.Redirect(w, r, "/admin?error=Write Error", http.StatusSeeOther); return }; defer out.Close(); _, err = io.Copy(out, file); initDBConnection(); if err != nil { http.Redirect(w, r, "/admin?error=Restore Failed", http.StatusSeeOther); return }; http.Redirect(w, r, "/admin?success=Restored!", http.StatusSeeOther) } }
