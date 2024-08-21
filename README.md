
# pfsms #

## Introduction ##
PFSMS is a tool for sending multiple SMS messages using a Mobile Phone connected via USB to the computer.
Since I myself use a Samsung Galaxy Phone, it is specifically configured for this phohne, but it will probably
work with other Phones supporting 3GPP AT-command set specifictions.
To enable AT-commands on my Samsung Galaxy phone, I need to 
(https://www.samsung.com/uk/support/mobile-devices/how-do-i-turn-on-the-developer-options-menu-on-my-samsung-galaxy-device/):
1 Go to "Settings"
2 Tap "About device" or "About phone"
3 Tap “Software information”
4 Tap “Build number” seven times. ...
5 Enter your pattern, PIN or password to enable the Developer options menu.
6 The "Developer options" menu will now appear in your Settings menu.
7 Then Scroll down amont options until "Error handling" and below there somewhere you will see:
"3GPP AT-commands". Turn this option ON !

## SMS messaging ##
To send an sms, just go to Message tab and either paste desired phonenumbers separated by commas (,) or return.
Or you can choose a group. A group is just a bunch of saved numbers with a groupname.
To add a group, specifya groupname and click "Save Group".
To make a copy of a group, select the group, change the groupname and click "Save Group".
Groups are saved into the Groups database, see below. Groups can also be exported and imported.
Phonenumbers are fixed before saved into a group.
When writing a message, you can insert firstname and lastname, by writing "<<fname>>" and "<<lname>>".

## Customers handling ##
It should be simple to handle, so it is not intended to be some kind of sms customer management system.
Therefore, all customer handling is recommended to be done using "Import/Export Customers" via a textfile.
All files are tab-separated textfiles and a sample can easily be obtained by exporting Customers.
The idea is to maintain the database of customers using some other solution (Google Sheets in my case).
The Customers database contains fields: id, phone, firstname, lastname, note.
When importing, all customers with incorrect or missing phonenumber will be skipped.
Only unique phonenumbers are possible to have in the database.

## Phonenumber fixing ##
All phones are checked and fixed before usage. Just digits are kept.
All leading "0" (just 1 zero) are replaced with "00", plus country code specified in settings.
This means all phonenumbers, when used to send sms will look like (for Sweden): "0046736290839"

## Messages ##
Messages sent are saved into as history into the database. History can be exported.
There is no real limit on sms messages. They work with special characters "åäöÅÄÖ".
It is not possible to send images in sms (MMS).

## Logging ##
pfsms.log - a logfile, where info on errors, success and other stuff happening in Program is saved
history.txt - a textfile whereto current history can be exported
customes.txt - a textfile for import/export of customers.
groups.txt - a textfile for import/export of groups.
pfsms.db - a simple SQL database with tables tblHistory, tblCustomers, tblGroups and tblHashtable

## Email ##
Email password is saved into preferences using a hashtable in database. It should be a little safer so...
The point of Email is to leave the system to get emails from an specific email account and then 
automatically send sms to specified phone(s) with specified message.
How to implement this I have still not decided.

## Settings ##
Specify PhoneNumber, Country, Model and Modem port and click to send a test sms to yoour own phone.
Then some import, export and opening solutions. Finally showing the latest part of pfsms.log.

## About
It is possible to enable DEBUG. This means that instead of the standard User directory, all files will be
saved into the folder "pfsms" within the same folder as the application. Practical for debuggning?!
Some info on database & memory & appearance is also available here!

For more information, please contact, with a smile
Peter Freund
peter.freund@prifre.com
