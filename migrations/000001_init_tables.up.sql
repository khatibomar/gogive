/*==============================================================*/
/* DBMS name:      PostgreSQL 14.1                              */
/* DBMS name:      Generated with power designer                */
/* Created on:     1/24/2022 12:40:00 AM                        */
/*==============================================================*/

/*==============================================================*/
/* Table: LOCATIONS                                             */
/*==============================================================*/
create table LOCATIONS (
   PCODE                TEXT                 not null,
   LOCATION_NAME_EN     TEXT                 UNIQUE not null,
   LOCATION_NAME_AR     TEXT                 UNIQUE not null,
   LATITUDE             FLOAT8               not null,
   LONGITUDE            FLOAT8               not null,
   GOVERNORATE_EN       TEXT                 not null,
   GOVERNORATE_AR       TEXT                 not null,
   VERSION              INT4                 not null default 1
      constraint CKC_VERSION_LOCATION check (VERSION >= 1),
   constraint PK_LOCATIONS primary key (PCODE)
);

/*==============================================================*/
/* Table: ROLES                                                 */
/*==============================================================*/
create table ROLES (
   ROLE_ID              Bigserial            not null,
   ROLE_NAME            TEXT                 UNIQUE not null,
   ROLE_DESCRIPTION     TEXT                 not null,
   constraint PK_ROLES primary key (ROLE_ID)
);

/*==============================================================*/
/* Table: USERS                                                 */
/*==============================================================*/
create table USERS (
   USER_ID              Bigserial            not null,
   ROLE_ID              bigint               not null,
   PCODE                TEXT                 not null,
   CREATED_AT           TIMESTAMP(0) with time zone not null DEFAULT NOW(),
   ACTIVATED            BOOL                 not null,
   PHOTO_URL            TEXT                 null,
   FIRSTNAME            TEXT                 not null,
   LASTNAME             TEXT                 not null,
   PHONE                TEXT                 null,
   EMAIL                Citext               UNIQUE not null,
   PASSWORD_HASH        bytea                not null,
   VERSION              INT4                 not null default 1
      constraint CKC_VERSION_USERS check (VERSION >= 1),
   constraint PK_USERS primary key (USER_ID),
   constraint FK_USERS_LIVES_IN_LOCATION foreign key (PCODE)
      references LOCATIONS (PCODE)
      on delete restrict on update restrict,
   constraint FK_USERS_HAVE_ROLE_ROLES foreign key (ROLE_ID)
      references ROLES (ROLE_ID)
      on delete restrict on update restrict
);

/*==============================================================*/
/* Table: BANS                                                  */
/*==============================================================*/
create table BANS (
   BANNED_BY_ID         bigserial            not null,
   USER_ID              bigint               not null,
   EMAIL                Citext               UNIQUE not null,
   BAN_REASON           TEXT                 not null,
   BAN_EXPIRY           TIMESTAMP with time zone            not null,
   VERSION              INT4                 not null default 1
      constraint CKC_VERSION_BANS check (VERSION >= 1),
   constraint PK_BANS primary key (BANNED_BY_ID),
   constraint FK_BANS_BAN_USERS foreign key (USER_ID)
      references USERS (USER_ID)
      on delete restrict on update restrict
);

/*==============================================================*/
/* Table: CATEGORIES                                            */
/*==============================================================*/
create table CATEGORIES (
   CATEGORY_ID          Bigserial            not null,
   CATEGORY_NAME        TEXT                 UNIQUE not null,
   constraint PK_CATEGORIES primary key (CATEGORY_ID)
);

/*==============================================================*/
/* Table: PERMISSIONS                                           */
/*==============================================================*/
create table PERMISSIONS (
   PERM_ID              Bigserial            not null,
   CODE                 TEXT                 UNIQUE not null,
   constraint PK_PERMISSIONS primary key (PERM_ID)
);

/*==============================================================*/
/* Table: HAVE_PERMISSION                                       */
/*==============================================================*/
create table HAVE_PERMISSION (
   PERM_ID              bigint               not null,
   ROLE_ID              bigint               not null,
   constraint PK_HAVE_PERMISSION primary key (PERM_ID, ROLE_ID),
   constraint FK_HAVE_PER_HAVE_PERM_PERMISSI foreign key (PERM_ID)
      references PERMISSIONS (PERM_ID)
      on delete restrict on update restrict,
   constraint FK_HAVE_PER_HAVE_PERM_ROLES foreign key (ROLE_ID)
      references ROLES (ROLE_ID)
      on delete restrict on update restrict
);

/*==============================================================*/
/* Table: ITEMS                                                 */
/*==============================================================*/
create table ITEMS (
   ITEM_ID              Bigserial            not null,
   PCODE                TEXT                 not null,
   USER_ID              bigint               not null,
   CATEGORY_ID          bigint               not null,
   CREATED_AT           TIMESTAMP(0) with time zone not null DEFAULT NOW(),
   VERSION              INT4                 not null default 1
      constraint CKC_VERSION_ITEMS check (VERSION >= 1),
   PHOTO_URL            TEXT                 null,
   constraint PK_ITEMS primary key (ITEM_ID),
   constraint FK_ITEMS_CREATE_USERS foreign key (USER_ID)
      references USERS (USER_ID)
      on delete restrict on update restrict,
   constraint FK_ITEMS_EXIST_IN_LOCATION foreign key (PCODE)
      references LOCATIONS (PCODE)
      on delete restrict on update restrict,
   constraint FK_ITEMS_CONTAINS_CATEGORI foreign key (CATEGORY_ID)
      references CATEGORIES (CATEGORY_ID)
      on delete restrict on update restrict
);

/*==============================================================*/
/* Table: TOKENS                                                */
/*==============================================================*/
create table TOKENS (
   HASH                 bytea                not null,
   USER_ID              bigint               UNIQUE not null,
   EXPIRY               TIMESTAMP			 not null,
   SCOPE                TEXT                 not null,
   constraint PK_TOKENS primary key (HASH),
   constraint FK_TOKENS_HAVE_TOKE_USERS foreign key (USER_ID)
      references USERS (USER_ID)
      on delete restrict on update restrict
);

/*==============================================================*/
/* Table: Initial Insertions									*/
/*==============================================================*/
INSERT INTO roles(role_name, role_description) VALUES
	('user', 'for normal user'),
	('admin','for admins'),
	('analytic','for analytics');
