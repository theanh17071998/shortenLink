CREATE KEYSPACE shortenurl
  WITH REPLICATION = { 
   'class' : 'SimpleStrategy', 
   'replication_factor' : 1 
  };

USE shortenurl;

CREATE TABLE url(
	shortcode text,
	longurl text,
	PRIMARY KEY(shortcode)
);

