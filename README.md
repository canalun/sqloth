
# sqloth
It is an offline SQL dummy data generator!!

*"Wanna test using dummy SQL data? You'll like SQLOTH...!"*
## 🎉 Features 🎉
- completely offline, which means you can use confidential schema
- automatically analyze foreign key dependencies and generate data along with them
- (not yet) ~~fast calculation, 1M records for XX secs!~~
  - (currently) 1000 records for around 100 columns are the limit...!
- (not yet) ~~variable formats for random data generation. you can set prefix, suffix and randomize methods(e.g. uuid)!~~
  - (currently) generate perfectly random data. you cannot set prefix, suffix or randomize methods...!
## 📦 Install 📦
Please download the binary. That's all!!
## 💻 Usage 💻
Please run the below.

```./sqloth -f ./path/to/your/schema.sql -n [the # of records you want]```

Here is an example of input and output.

```
$ sqloth -f ./path/to/your/schema.sql -n 10 > dummy.sql

$ cat dummy.sql

SET foreign_key_checks = 0;

INSERT INTO customer(`created_at`, `name`, `material`)
VALUES ('1982-02-12 12:22:27','Lhras20e...r7U3','{"json":"7647947524"}'),
...
('2021-11-05 11:32:13','aioI...I5t','{"json":"8493280504"}'),
('2004-05-11 00:57:27','86MI...PVn','{"json":"7486664121"}');

INSERT INTO product(`name`, `owner`, `description`, `stock`, `sale_day`)
VALUES ('Eq...fW','Lhr...U3','gILE...FDvK','0','2015-10-30 05:21:22'),
...
('SQU..62v','waN...Imm','kwL...gh8','1','2010-01-30 14:51:37'),
('ceJ...3xl','KvR...1Nm','NN4...vky','0','2022-03-08 05:43:08');

SET foreign_key_checks = 1;
```
## 🌟 Contribution 🌟
- Let's be creative and collaborative👶
- Please read [CONTRIBUTING.md](https://github.com/canalun/sqloth/blob/main/CONTRIBUTING.md) for the details😉
