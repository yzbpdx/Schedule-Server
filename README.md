# Schedule-Server

## mysql

```{.line-numbers}
create table students {
    studentId int primary key,
    studentName varchar(255),
    class varchar(255),
    classmates JSON,
    spareTime JSON,
    lesson JSON,
};
```