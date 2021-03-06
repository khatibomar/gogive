/*==============================================================*/
/* Index: BANS_PK                                               */
/*==============================================================*/
create unique index IF NOT EXISTS BANS_PK on BANS (
ID
);

/*==============================================================*/
/* Index: BAN_FK                                                */
/*==============================================================*/
create  index IF NOT EXISTS BAN_FK on BANS (
USER_ID
);

/*==============================================================*/
/* Index: CATEGORIES_PK                                         */
/*==============================================================*/
create unique index IF NOT EXISTS CATEGORIES_PK on CATEGORIES (
ID
);

/*==============================================================*/
/* Index: ITEMS_PK                                              */
/*==============================================================*/
create unique index IF NOT EXISTS ITEMS_PK on ITEMS (
ID
);

/*==============================================================*/
/* Index: CREATE_FK                                             */
/*==============================================================*/
create  index IF NOT EXISTS CREATE_FK on ITEMS (
USER_ID
);

/*==============================================================*/
/* Index: EXIST_IN_FK                                           */
/*==============================================================*/
create  index IF NOT EXISTS EXIST_IN_FK on ITEMS (
PCODE
);

/*==============================================================*/
/* Index: BELONGS_FK                                            */
/*==============================================================*/
create  index IF NOT EXISTS BELONGS_FK on ITEMS (
CATEGORY_ID
);

/*==============================================================*/
/* Index: LOCATIONS_PK                                          */
/*==============================================================*/
create unique index IF NOT EXISTS LOCATIONS_PK on LOCATIONS (
PCODE
);

/*==============================================================*/
/* Index: TOKENS_PK                                             */
/*==============================================================*/
create unique index IF NOT EXISTS TOKENS_PK on TOKENS (
HASH
);

/*==============================================================*/
/* Index: HAVE_TOKEN_FK                                         */
/*==============================================================*/
create  index IF NOT EXISTS HAVE_TOKEN_FK on TOKENS (
USER_ID
);

/*==============================================================*/
/* Index: USERS_PK                                              */
/*==============================================================*/
create unique index IF NOT EXISTS USERS_PK on USERS (
ID
);

/*==============================================================*/
/* Index: LIVES_IN_FK                                           */
/*==============================================================*/
create index IF NOT EXISTS LIVES_IN_FK on USERS (
PCODE
);
