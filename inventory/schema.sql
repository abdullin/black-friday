PRAGMA journal_mode=WAL;

PRAGMA foreign_keys = ON;

create table dummy (id integer primary key autoincrement);


CREATE TABLE Locations (

    Id INTEGER PRIMARY KEY,
    -- can be null for root
    Parent INTEGER,
    Name TEXT NOT NULL UNIQUE,
    FOREIGN KEY(Parent) REFERENCES Locations(Id)
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
    ('Products', 0);



