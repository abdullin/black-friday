PRAGMA journal_mode=WAL;

PRAGMA foreign_keys = ON;

create table dummy (id integer primary key autoincrement);

CREATE TABLE Warehouses(
    Id INTEGER PRIMARY KEY,
    Name TEXT NOT NULL UNIQUE
);

CREATE TABLE Locations (
    Id INTEGER PRIMARY KEY,
    Warehouse INTEGER NOT NULL,
    Name TEXT NOT NULL UNIQUE,
    FOREIGN KEY(Warehouse) REFERENCES Warehouses(Id)
);


CREATE TABLE Products (
    Id INTEGER PRIMARY KEY,
    Sku TEXT NOT NULL UNIQUE
);

CREATE TABLE Inventory (
    Location INTEGER NOT NULL,
    Product INTEGER NOT NULL,
    OnHand INTEGER NOT NULL,
    FOREIGN KEY(Location) REFERENCES Locations(Id),
    FOREIGN KEY(Product) REFERENCES Products(Id)
);

INSERT INTO sqlite_sequence (name, seq) VALUES
    ('Locations', 0),
    ('Products', 0),
    ('Warehouses', 0);



