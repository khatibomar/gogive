/*==============================================================*/
/* Index: BANS_PK                                               */
/*==============================================================*/
create unique index IF NOT EXISTS BANS_PK on BANS (
BANNED_BY_ID
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
CATEGORY_ID
);

/*==============================================================*/
/* Index: HAVE_PERMISSION_FK                                    */
/*==============================================================*/
create  index IF NOT EXISTS HAVE_PERMISSION_FK on HAVE_PERMISSION (
PERM_ID
);

/*==============================================================*/
/* Index: HAVE_PERMISSION2_FK                                   */
/*==============================================================*/
create  index IF NOT EXISTS HAVE_PERMISSION2_FK on HAVE_PERMISSION (
ROLE_ID
);

/*==============================================================*/
/* Index: ITEMS_PK                                              */
/*==============================================================*/
create unique index IF NOT EXISTS ITEMS_PK on ITEMS (
ITEM_ID
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
/* Index: CONTAINS_FK                                           */
/*==============================================================*/
create  index IF NOT EXISTS CONTAINS_FK on ITEMS (
CATEGORY_ID
);

/*==============================================================*/
/* Index: LOCATIONS_PK                                          */
/*==============================================================*/
create unique index IF NOT EXISTS LOCATIONS_PK on LOCATIONS (
PCODE
);

/*==============================================================*/
/* Index: PERMISSIONS_PK                                        */
/*==============================================================*/
create unique index IF NOT EXISTS PERMISSIONS_PK on PERMISSIONS (
PERM_ID
);

/*==============================================================*/
/* Index: ROLES_PK                                              */
/*==============================================================*/
create unique index IF NOT EXISTS ROLES_PK on ROLES (
ROLE_ID
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
USER_ID
);

/*==============================================================*/
/* Index: LIVES_IN_FK                                           */
/*==============================================================*/
create  index IF NOT EXISTS LIVES_IN_FK on USERS (
PCODE
);

/*==============================================================*/
/* Index: HAVE_ROLE_FK                                          */
/*==============================================================*/
create  index IF NOT EXISTS HAVE_ROLE_FK on USERS (
ROLE_ID
);
