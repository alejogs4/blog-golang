CREATE DATABASE blog

CREATE TABLE person (
  id VARCHAR(140),
  firstname VARCHAR(200) NOT NULL,
  lastname VARCHAR(200) NOT NULL,
  email VARCHAR(300) NOT NULL,
  email_verified Boolean NOT NULL DEFAULT false,
  password VARCHAR(300) NOT NULL,
  CONSTRAINT person_pk PRIMARY KEY(id),
  CONSTRAINT person_email_unique UNIQUE(email)
);

CREATE TABLE post (
  id VARCHAR(140),
  person_id VARCHAR(140) NOT NULL,
  title VARCHAR(300) NOT NULL,
  content TEXT NOT NULL,
  picture VARCHAR(400) NOT NULL,
  CONSTRAINT post_pk PRIMARY KEY(id),
  CONSTRAINT post_person_fk FOREIGN KEY(person_id) REFERENCES person(id)
);

create type comment_state as enum('ACTIVE', 'REMOVED');
CREATE TABLE comment (
  id VARCHAR(140),
  content TEXT NOT NULL,
  person_id VARCHAR(140),
  post_id VARCHAR(140),
  state comment_state NOT NULL,
  CONSTRAINT comment_pk PRIMARY KEY(id),
  CONSTRAINT comment_person_fk FOREIGN KEY(person_id) REFERENCES person(id),
  CONSTRAINT comment_post_fk FOREIGN KEY(post_id) REFERENCES post(id)
);


create type like_state as enum('ACTIVE', 'REMOVED');
create type like_type as enum('LIKE', 'DISLIKE');
CREATE TABLE post_like (
  id VARCHAR(140),
  person_id VARCHAR(140),
  post_id VARCHAR(140),
  state like_state NOT NULL,
  type like_type NOT NULL,
  CONSTRAINT like_pk PRIMARY KEY(id),
  CONSTRAINT like_person FOREIGN KEY(person_id) REFERENCES person(id),
  CONSTRAINT like_post FOREIGN KEY(post_id) REFERENCES post(id)
);

CREATE TABLE tag (
  id VARCHAR(140),
  content VARCHAR(120) NOT NULL,
  CONSTRAINT tag_pk PRIMARY KEY(id)
);

CREATE TABLE post_tag (
  tag_id VARCHAR(140),
  post_id VARCHAR(140),
  CONSTRAINT post_tag_pk PRIMARY KEY(tag_id, post_id),
  CONSTRAINT post_tag_tag FOREIGN KEY(tag_id) REFERENCES tag(id),
  CONSTRAINT post_tag_post FOREIGN KEY(post_id) REFERENCES post(id)
);