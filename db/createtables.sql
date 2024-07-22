CREATE TABLE tblMessages (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   nanostamp integer, tstamp VARCHAR(25), reference VARCHAR(100), message TEXT);

CREATE TABLE tblCustomers (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   expnote integer, phone VARCHAR(20), firstname VARCHAR(100), lastname VARCHAR(100), 
   indate VARCHAR(10), outdate VARCHAR(10), note TEXT);

CREATE TABLE tblHistory (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   nanostamp integer, tstamp VARCHAR(20), 
   customerid integer, phone VARCHAR(20),
   messageid integer, reference VARCHAR(100), message TEXT);

CREATE TABLE tblGroups (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   groupid integer, customerid integer);

CREATE TABLE tblGroupnames (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   groupname VARCHAR(100));

CREATE TABLE tblHashtable (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   hash VARCHAR(100));
