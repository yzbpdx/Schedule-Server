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

insert into teachers values(8, 'Erin', '{"0": {"0": {}, "1": {}, "2": {}}, "1": {"0": {}, "1": {}, "2": {}}, "2": {"0": {}, "1": {}, "2": {}}, "3": {"0": {}, "1": {}, "2": {}}, "4": {"0": {}, "1": {}, "2": {}}, "5": {"0": {}, "1": {}, "2": {}}, "6": {"0": {}, "1": {}, "2": {}}}', 2, true);

insert into students values(7, '张颖严', '{"0": {"0": {}, "1": {}, "2": {}}, "1": {"0": {}, "1": {}, "2": {}}, "2": {"0": {}, "1": {}, "2": {}}, "3": {"0": {}, "1": {}, "2": {}}, "4": {"0": {}, "1": {}, "2": {}}, "5": {"0": {}, "1": {}, "2": {}}, "6": {"0": {}, "1": {}, "2": {}}}', true);
```