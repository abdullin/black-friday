
create table dummy (id integer primary key autoincrement);

CREATE TABLE Locations (
    Id INTEGER PRIMARY KEY,
    Name TEXT NOT NULL UNIQUE
);


CREATE TABLE Products (
    Id INTEGER PRIMARY KEY,
    Sku TEXT NOT NULL UNIQUE
);

CREATE TABLE Inventory (
    Location INTEGER NOT NULL,
    Product INTEGER NOT NULL,
    Quantity INTEGER NOT NULL,
    FOREIGN KEY(Location) REFERENCES Location(Id),
    FOREIGN KEY(Product) REFERENCES Products(Id)
);

INSERT INTO sqlite_sequence (name, seq) VALUES
    ('Locations', 0),
    ('Products', 0);



