CREATE TABLE Organisation
(
    id       INT PRIMARY KEY AUTO_INCREMENT,
    name     VARCHAR(256) NOT NULL,
    password VARCHAR(256) NOT NULL,
    UNIQUE KEY (name)
);

CREATE TABLE Mousetrap
(
    id        INT PRIMARY KEY AUTO_INCREMENT,
    name      VARCHAR(256) NOT NULL,
    org_id    INT          NOT NULL,
    status    BOOLEAN      NOT NULL,
    last_trig INT          NOT NULL,
    FOREIGN KEY (org_id) REFERENCES Organisation (id) ON DELETE CASCADE,
    UNIQUE KEY (name, org_id)
);