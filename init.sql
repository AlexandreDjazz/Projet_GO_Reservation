CREATE TABLE SALLES(
   id_salle INT AUTO_INCREMENT,
   nom VARCHAR(50) NOT NULL,
   place INT NOT NULL,
   PRIMARY KEY(id_salle)
);

CREATE TABLE ETATS(
   id_etat INT AUTO_INCREMENT,
   nom_etat VARCHAR(50) NOT NULL,
   PRIMARY KEY(id_etat)
);

CREATE TABLE RESERVATIONS(
   id_reservation INT AUTO_INCREMENT,
   horaire_start DATETIME NOT NULL,
   horaire_end DATETIME NOT NULL,
   id_etat INT NOT NULL,
   PRIMARY KEY(id_reservation),
   FOREIGN KEY(id_etat) REFERENCES ETATS(id_etat)
);

CREATE TABLE RESERVER(
   id_salle INT,
   id_reservation INT,
   PRIMARY KEY(id_salle, id_reservation),
   FOREIGN KEY(id_salle) REFERENCES SALLES(id_salle),
   FOREIGN KEY(id_reservation) REFERENCES RESERVATIONS(id_reservation)
);

INSERT INTO SALLES (id_salle, nom, place) VALUES
(1, 'Salle A', 50),
(2, 'Salle B', 40),
(3, 'Salle C', 60),
(4, 'Salle D', 50),
(5, 'Salle E', 40),
(6, 'Salle F', 60),
(7, 'Salle G', 50),
(8, 'Salle H', 40),
(9, 'Salle I', 60),
(10, 'Salle J', 60),
(11, 'Salle K', 50),
(12, 'Salle L', 40),
(13, 'Salle 13', 69);

INSERT INTO ETATS (id_etat, nom_etat) VALUES
(1, 'En cours'),
(2, 'Terminé'),
(3, 'Annulé'),
(4,'A venir');

INSERT INTO RESERVATIONS (id_reservation, horaire_start, horaire_end, id_etat) VALUES
(1, '2024-04-13 10:00:00','2024-04-13 11:00:00', 1),
(2, '2024-04-14 15:30:00','2024-04-14 17:00:00', 2),
(3, '2024-04-15 15:30:00','2024-04-15 17:30:00', 3),
(4, '2024-04-15 12:00:00','2024-04-15 14:00:00', 1);

INSERT INTO RESERVER (id_salle, id_reservation) VALUES
(1, 1),
(2, 2),
(1, 3),
(3, 4);