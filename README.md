What is Kasi? Content:

Kasi: The Privacy-First, Self-Hosted Budget Manager

Kasi is a lightweight, open-source expense tracking solution designed for individuals, freelancers, and project managers who value data privacy. Unlike traditional finance apps that store your financial data on third-party servers, Kasi is designed to be Self-Hosted. This means your data stays on your own server or device—100% under your control.

Why Choose Kasi?

🔒 Privacy-First: No tracking, no ads, and no selling of your data.

🚀 Ultra-Lightweight: Built with Go (Golang) and SQLite, it runs efficiently on minimal resources (requires less than 20MB RAM).

📱 Native-Like Experience: Works as a Progressive Web App (PWA) on both iOS and Android. Install it directly to your home screen without an App Store.

🌍 Global Ready: Supports multi-currency (USD, EUR, GBP, LKR, etc.) and multiple languages.

📊 Smart Reporting: Generate printable PDF reports for monthly expenses or specific events (e.g., Weddings, Trips, Construction projects).

Whether you are managing a family budget, a freelance project, or a wedding event, Kasi adapts to your needs with a simple, clean interface.

📖 Kasi User Manual: Getting Started
Welcome to Kasi! This guide will help you set up and manage your expenses efficiently.

🚀 Phase 1: Installation (For Self-Hosters)
Kasi is distributed via Docker, making it incredibly easy to install on any server (VPS, Raspberry Pi, or Laptop).

Run this command in your terminal:

Bash

docker run -d \
  -p 8080:8080 \
  -v ./kasi-data:/app/data \
  -e SESSION_SECRET="ReplaceWithAStrongPassword" \
  --name kasi \
  --restart unless-stopped \
  amilakothalawalasolo/kasi:latest
Port: The app will run on port 8080.

Data: All data is saved in the ./kasi-data folder on your machine.

Security: Change ReplaceWithAStrongPassword to something unique.

📱 Phase 2: Installing on Mobile (PWA)
Kasi functions exactly like a native app. You don't need the App Store or Play Store.

🤖 For Android Users
Open Kasi in Chrome.

Go to Settings (⚙️ icon in the app).

Tap the "Install on Android" button.

Confirm to add it to your Home Screen.

🍏 For iPhone (iOS) Users
Open Kasi in Safari.

Tap the Share Button (Square with an arrow up).

Scroll down and select "Add to Home Screen".

Tap Add. The app will now appear on your home screen and open in full-screen mode.

⚙️ Phase 3: Initial Setup
Once installed, follow these steps to personalize your experience:

Login:

Default Admin Username: admin

Default Password: admin123

(Note: Change these immediately after logging in!)

Go to Settings:

Click the Gear Icon (⚙️) on the dashboard.

Update Profile:

Display Name: Your name.

Project Name: Give your budget a name (e.g., "Family Expenses", "Wedding 2026", "Trip to Bali"). This name will appear on the Dashboard and Reports for all users.

Currency: Select your preferred currency (USD, EUR, LKR, etc.).

Password: Set a secure new password.

💸 Phase 4: Managing Expenses
The Dashboard is designed for speed.

Adding an Expense
Item Name: What did you buy? (e.g., "Dinner", "Cement", "Fuel").

Category: Choose from the list (Food, Transport, etc.).

Qty/Unit: (Optional) Useful for construction or shopping lists (e.g., 5 kg).

Price: Enter the amount.

Type:

Private: Only visible to you.

Common: Visible to all users in the project (e.g., Admin and Spouse).

Click Add.

Editing/Deleting
Tap the Pencil Icon next to an item to edit the amount or name.

Tap the Bin Icon to delete an entry (Admins can delete any entry; Users can only delete their own).

📊 Phase 5: Reports & Analytics
Need to see where your money went?

On the Dashboard, scroll down to the "Filter by Date" section.

Select a Start Date and End Date.

Click Search to filter the list.

Click the green "📄 Generate Report" button.

This opens a clean, printable view.

You can save this as a PDF or print it directly for your records.

👑 Phase 6: Admin Features
If you are the Admin (the owner of the server), you have extra powers:

Backup Data:

Go to Admin Panel (Crown Icon 👑).

Click "Download Backup" to save your .db file locally.

Tip: Do this weekly.

Restore Data:

If you move servers, simply upload your backup file in the "Restore" section.

Manage Users:

Create new accounts for family members or team members.

Reset their passwords if they forget them.

❤️ Support the Project
Kasi is an open-source project developed by Amila Kothalawala. If you find it useful, you can support the development via crypto donations
