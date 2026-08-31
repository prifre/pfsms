# PFSMS

PFSMS is a desktop tool for sending bulk SMS messages using a mobile phone connected to your computer via USB.

It is tailored and optimized for modern Samsung Galaxy devices (e.g., S24 Ultra), but it works with any phone supporting the standard **3GPP AT-command set**.

---

## 📱 Device Setup (Samsung Galaxy)

To allow PFSMS to communicate with your Samsung phone via serial AT-commands, follow these steps:

1. Go to **Settings** > **About phone** > **Software information**.
2. Tap **Build number** 7 times (enter your PIN/Password if prompted) to enable **Developer options**.
3. Go back to **Settings** > **Developer options**.
4. Enable **USB Debugging**.
5. Set **Default USB configuration** to **Transferring files**.
6. Scroll down to **3GPP AT-commands** and turn it **ON**.
7. Disable **Auto Blocker**: Go to **Settings** > **Security and privacy** > **Auto Blocker** and turn it **OFF** (this feature blocks serial USB communication).
8. Ensure you have installed the official **Samsung USB Drivers** on your PC ([Download Here](https://developer.samsung.com/android-usb-driver)).

---

## ✨ Features & Usage

### 📩 SMS Messaging
- **Direct Sending:** Go to the *Messages* tab and paste phone numbers separated by commas or line breaks.
- **Group Management:** Enter a group name and click **Save Group**. To duplicate a group, select it, change the name, and save.
- **Dynamic Tags:** Insert dynamic variables like `<<Fname>>` and `<<Lname>>` in your text to personalize messages.
- **Special Characters:** Full support for UTF-8 and Swedish characters (`å`, `ä`, `ö`, etc.).

### 👥 Customer Management
Customer data is kept minimal for performance. Full management is best handled via **Import/Export** using tab-separated text files (`.txt` / `.csv`) or external tools like Google Sheets.
- **Fields:** `id`, `phone`, `firstname`, `lastname`, `note`
- **Validation:** Automatic duplicate removal (only unique phone numbers are allowed) and invalid number filtering.

### 📜 Message History
All sent messages are stored in a local SQLite database. History can be searched or exported to text files at any time.

### ⚙️ Settings & Debugging
- Configure target serial modem port, country code, and phone model.
- Includes a **Send Test SMS** function to verify modem connectivity.
- **Debug Mode:** Enable debug mode in the *About* tab to store all database and log files in the application folder instead of the user's home directory.

---

## 📁 File Structure

| File / Database | Description |
| :--- | :--- |
| `pfsms.db` | Local SQLite database (`tblHistory`, `tblCustomers`, `tblGroups`, `tblHashtable`, `tblQueue`) |
| `pfsms.log` | Application runtime log file for troubleshooting |
| `customers.txt` | Import/Export file for customer data (tab-separated) |
| `groups.txt` | Import/Export file for contact groups |
| `history.txt` | Exported SMS dispatch history |

---

## 📧 Automated Email-to-SMS (Planned)
The application includes an encrypted hashtable preference store for email credentials. Future releases aim to poll a specific email account and automatically forward incoming emails as SMS messages to specified targets.

---

## 👤 Author & Support

**Peter Freund**  
📧 Email: [peter.freund@prifre.com](mailto:peter.freund@prifre.com)  
🌐 Website: [https://sms.prifre.com](https://sms.prifre.com)