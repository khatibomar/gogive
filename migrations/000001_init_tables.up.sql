/*==============================================================*/
/* Table: BANS                                                  */
/*==============================================================*/
create table IF NOT EXISTS BANS (
   ID					bigserial            not null,
   USER_ID              bigint				 null,
   EMAIL                Citext               UNIQUE not null,
   BAN_REASON           TEXT                 not null,
   BAN_EXPIRY           DATE                 not null,
   VERSION              INT4                 not null default 1
      constraint CKC_VERSION_BANS check (VERSION >= 1),
   constraint PK_BANS primary key (ID)
);

/*==============================================================*/
/* Table: CATEGORIES                                            */
/*==============================================================*/
create table IF NOT EXISTS CATEGORIES (
   ID Bigserial            not null,
   CATEGORY_NAME        TEXT                 UNIQUE not null
      constraint CKC_CATEGORY_NAME_CATEGORI check (CATEGORY_NAME = lower(CATEGORY_NAME)),
   constraint PK_CATEGORIES primary key (ID)
);

/*==============================================================*/
/* Table: ITEMS                                                 */
/*==============================================================*/
create table IF NOT EXISTS ITEMS (
   ID                   Bigserial            not null,
   PCODE                TEXT                 null,
   USER_ID              bigint               null,
   CATEGORY_ID          bigint               null,
   CREATED_AT           TIMESTAMP(0) with time zone not null DEFAULT NOW(),
   NAME                 TEXT                 not null
      constraint CKC_NAME_ITEMS check (NAME = lower(NAME)),
   QUANTITY             INT4                 not null default 1,
   IMAGE_URL            TEXT                 null,
   VERSION              INT4                 not null default 1
      constraint CKC_VERSION_ITEMS check (VERSION >= 1),
   constraint PK_ITEMS primary key (ID)
);

/*==============================================================*/
/* Table: LOCATIONS                                             */
/*==============================================================*/
create table IF NOT EXISTS LOCATIONS (
   PCODE                TEXT                 not null,
   LOCATION_NAME_EN     TEXT                 unique not null,
   LOCATION_NAME_AR     TEXT                 unique not null,
   LATITUDE             FLOAT8               not null,
   LONGITUDE            FLOAT8               not null,
   GOVERNORATE_EN       TEXT                 not null,
   GOVERNORATE_AR       TEXT                 not null,
   VERSION              INT4                 not null default 1
      constraint CKC_VERSION_LOCATION check (VERSION >= 1),
   constraint PK_LOCATIONS primary key (PCODE)
);

/*==============================================================*/
/* Table: TOKENS                                                */
/*==============================================================*/
create table IF NOT EXISTS TOKENS (
   HASH                 bytea                not null,
   USER_ID              bigint               null,
   EXPIRY               DATE                 not null,
   SCOPE                TEXT                 not null,
   constraint PK_TOKENS primary key (HASH)
);

/*==============================================================*/
/* Table: USERS                                                 */
/*==============================================================*/
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'userclass') THEN
		create type userclass as enum('user','analytic','admin');
    END IF;
END$$;

create table IF NOT EXISTS USERS (
   ID 					Bigserial            not null,
   PCODE                TEXT                 null,
   CREATED_AT           TIMESTAMP(0) with time zone not null DEFAULT NOW(),
   ACTIVATED            BOOL                 not null,
   IMAGE_URL            TEXT                 null,
   FIRSTNAME            TEXT                 not null,
   LASTNAME             TEXT                 not null,
   PHONE                TEXT                 null,
   EMAIL                Citext               UNIQUE not null,
   PASSWORD_HASH        bytea                not null,
   USERTYPE             userclass            not null default 'user',
   IS_ANALYTIC          BOOL                 not null default false, 
   VERSION              INT4                 not null default 1
      constraint CKC_VERSION_USERS check (VERSION >= 1),
   constraint PK_USERS primary key (ID)
);

/*==============================================================*/
/* Table: FOREIGN KEYS                                          */
/*==============================================================*/
alter table BANS
   add constraint FK_BANS_BAN_USERS foreign key (USER_ID)
      references USERS (ID)
      on delete cascade on update cascade;

alter table ITEMS
   add constraint FK_ITEMS_BELONGS_CATEGORI foreign key (CATEGORY_ID)
      references CATEGORIES (ID)
      on delete cascade on update cascade;

alter table ITEMS
   add constraint FK_ITEMS_CREATE_USERS foreign key (USER_ID)
      references USERS (ID)
	  on delete cascade on update cascade;

alter table ITEMS
   add constraint FK_ITEMS_EXIST_IN_LOCATION foreign key (PCODE)
      references LOCATIONS (PCODE)
	  on delete set null on update cascade;

alter table TOKENS
   add constraint FK_TOKENS_HAVE_TOKE_USERS foreign key (USER_ID)
      references USERS (ID)
	  on delete cascade on update cascade;

alter table USERS
   add constraint FK_USERS_LIVES_IN_LOCATION foreign key (PCODE)
      references LOCATIONS (PCODE)
      on delete set null on update cascade;
