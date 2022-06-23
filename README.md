# sqloth
offline sql faker generator
# features
- completely offline, which means you can use confidential schema!
- automatically analyze foreign key dependency and generate data along with it!
- (not yet) fast calculation, 1M records for XX secs!
  - (currently) 1000 records for around 100 columns are the limit...!
- (not yet) variable formats for random data generation. you can set prefix, suffix and randomize methods(e.g. uuid)!
  - (currently) perfectly random. you cannot set prefix, suffix and randomize methods...!
# usage
please download the binary and run the below!!

```./sqloth -f ./path/to/your/schema.sql -n [the # of records you want]```

the below is the example input and output.

```
$ sqloth -f ./path/to/your/schema.sql -n 10 > dummy.sql

$ cat dummy.sql

SET foreign_key_checks = 0;

INSERT INTO customer(`created_at`, `name`, `material`) VALUES ('1982-02-12 12:22:27','Lhras20e...r7U3','{"json":"7647947524"}'),
...
('2021-11-05 11:32:13','aioI...I5t','{"json":"8493280504"}'),
('2004-05-11 00:57:27','86MI...PVn','{"json":"7486664121"}');

INSERT INTO product(`name`, `owner`, `description`, `stock`, `sale_day`) VALUES ('Eq...fW','Lhr...U3','gILE...FDvK','0','2015-10-30 05:21:22'),
...
('SQU..62v','waN...Imm','kwL...gh8','1','2010-01-30 14:51:37'),
('ceJ...3xl','KvR...1Nm','NN4...vky','0','2022-03-08 05:43:08');

SET foreign_key_checks = 1;
```
# contribution
- be creative and collaborative!!
- please read contribution.md for the details!!
