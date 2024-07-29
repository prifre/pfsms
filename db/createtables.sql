CREATE TABLE tblMessages (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   tstamp VARCHAR(25),
   reference VARCHAR(100), message TEXT);

CREATE TABLE tblCustomers (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   phone VARCHAR(20), firstname VARCHAR(100), lastname VARCHAR(100), note TEXT);

CREATE TABLE tblGroups (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   groupname varchar(100), phone varchar(100));

CREATE TABLE tblHistory (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   tstamp VARCHAR(20), 
   phone VARCHAR(20), reference VARCHAR(100), message TEXT);

CREATE TABLE tblHashtable (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   hash VARCHAR(100));
