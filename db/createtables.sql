CREATE TABLE tblMessages (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   nanostamp integer, tstamp TEXT, messagetitle VARCHAR(100), message TEXT)

CREATE TABLE tblCustomers (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
   expnote integer, phone VARCHAR(20), firstname VARCHAR(100), lastname VARCHAR(100), 
   indate VARCHAR(10), outdate VARCHAR(10),note TEXT);
