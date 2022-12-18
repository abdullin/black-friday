PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL ;

PRAGMA foreign_keys = ON;
PRAGMA cache_size = -8000;

create table dummy (id integer primary key autoincrement);


CREATE TABLE Locations (
    Id INTEGER PRIMARY KEY,
    Parent INTEGER NOT NULL,
    Name TEXT NOT NULL,
    FOREIGN KEY(Parent) REFERENCES Locations(Id),
    UNIQUE (Name, Parent)
);

CREATE INDEX IDX_LOCATIONS_PARENT
    ON Locations (Parent, Id);


CREATE TABLE Products (
    Id INTEGER PRIMARY KEY,
    Sku TEXT NOT NULL UNIQUE
);

CREATE UNIQUE INDEX IDX_PRODUCTS_SKU
    ON Products(Sku);


CREATE TABLE Inventory (
    Location INTEGER NOT NULL,
    Product INTEGER NOT NULL,
    OnHand INTEGER NOT NULL,
    FOREIGN KEY(Location) REFERENCES Locations(Id),
    FOREIGN KEY(Product) REFERENCES Products(Id),
    primary key (Product, Location)
) WITHOUT ROWID;

CREATE INDEX IDX_INVENTORY_PRODUCT
    ON Inventory (Product);

CREATE INDEX IDX_INVENTORY_LOCATION
    ON Inventory (Location);

CREATE UNIQUE INDEX IDX_INVENTORY_PRODUCT_LOCATION
    ON Inventory (Location, Product);

CREATE TABLE Reservations (
    Id INTEGER PRIMARY KEY,
    Code TEXT NOT NULL UNIQUE
);

CREATE TABLE Reserves (
    Reservation INTEGER NOT NULL,
    Product INTEGER NOT NULL,
    Location INTEGER NOT NULL,
    Quantity INTEGER NOT NULL,
    FOREIGN KEY(Reservation) REFERENCES Reservations(Id),
    FOREIGN KEY(Product) REFERENCES Products(Id),
    FOREIGN KEY(Location) REFERENCES Locations(Id),
    primary key (Reservation, Product, Location)
) WITHOUT ROWID ;


CREATE INDEX IDX_RESERVES_PRODUCT_LOCATION
    ON Reserves (Location, Product);

CREATE INDEX IDX_RESERVES_RESERVATION
    ON Reserves (Reservation);

INSERT INTO sqlite_sequence (name, seq) VALUES ('Entity', 0);


INSERT INTO Locations(Id, Parent, Name) VALUES(0,0, "Root");


