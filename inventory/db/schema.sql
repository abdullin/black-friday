PRAGMA journal_mode=WAL;

PRAGMA foreign_keys = ON;

create table dummy (id integer primary key autoincrement);


CREATE TABLE Locations (
    Id INTEGER PRIMARY KEY,
    Parent INTEGER NOT NULL,
    Name TEXT NOT NULL,
    FOREIGN KEY(Parent) REFERENCES Locations(Id),
    UNIQUE (Name, Parent)
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

CREATE TABLE Reservations (
    Id INTEGER PRIMARY KEY,
    Code TEXT NOT NULL UNIQUE
);

CREATE TABLE Reserves (
    Reservation INTEGER NOT NULL,
    Product INTEGER NOT NULL,
    Quantity INTEGER NOT NULL,
    -- CAN be null for ROOT
    -- TODO: Introduce root location of zero?
    Location INTEGER NOT NULL,
    FOREIGN KEY(Reservation) REFERENCES Reservations(Id),
    FOREIGN KEY(Product) REFERENCES Products(Id),
    FOREIGN KEY(Location) REFERENCES Locations(Id)
);

INSERT INTO sqlite_sequence (name, seq) VALUES
    ('Locations', 0),
    ('Products', 0),
    ('Reservations', 0),
    ('Global', 0);

INSERT INTO Locations(Id, Parent, Name) VALUES(0,0, "Root");


