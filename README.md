# pfsms

PFSMS is a tool for sending multiple SMS messages using a Mobile Phone connected via USB to the computer.
Since I myself use a Samsung Galaxy Phone, it is specifically configured for this phohne, but it will probably
work with other Phones supporting 3GPP AT-command set specifictions.

SMS messaging
To send an sms, just go to Message tab and either paste desired phonenumbers separated by commas or choose a group.
If more than one phone is specified, you will also have to specify a group name. A group will then be created.
Groups are saved into the Groups database, see below.
When writing a message, you can insert firstname, lastname and note, by writing 
"<<fname>>", "<<lname>>" and "<<note>>".

Customers handling
It should be simple to handle, so it is not intended to be some kind of sms customer management system.
Therefore, all customer handling is recommended to be done using "Import/Export Customers" via a textfile.
All files are tab-separated textfiles and a sample can easily be obtained by first just exporting.
The idea is to maintain the database of customers using some other solution (Google Sheets in my case).
The Customers database contains fields: id, phone firstname, lastname, note.
When importing, all customers with incorrect or missing phonenumber will be skipped.
Only unique phonenumbers are possible to have in the database.

Phonenumber fixing
All phones are checked and fixed before usage. This means:
- all characters, except "+","0","1","2","3","4","5","6","7","8","9" are removed
- all "+" signs are replaced with "00" (where what follows is assumed to be correct country code)
- all leading "00" are treated as correct, complete phonenumber, that including country
- all leading "0" (just 1 zero) are replaced with "00", plus country code specified in settings.
This means all phonenumbers, when used to send sms will look like (for Sweden): "0046736290839"

Groups handling
It should be easy to send sms to just specific customers, based on Groups.
Groups can be maintained using "Import/Export Groups".
Additional groups can be created using by sending messages to multiple phonesnumbers.

Messages
There is no real limit on sms messages. They work with special characters "åäöÅÄÖ".
It is not possible to send images in sms (MMS).

Program was developed using Go and Fyne, originally for Windows.
It should be possible to adjust it for Linux and Mac fairly easily.

For more information, please contact, with a smile
Peter Freund
peter.freund@prifre.com
