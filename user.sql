create table users (
    username varchar(128) primary key,
    pwhash varchar(64) not null
)
