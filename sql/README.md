# Legacy data import from report interface

# Edits

```bash
mysql --defaults-file="${HOME}"/replica.my.cnf -h tools-db s52585__cb -s -r -e 'select concat("INSERT INTO edit VALUES (", new_id, ", 1, 0, ", if(status=7,1,0), ");") from reports inner join vandalism on id=revertid where status in (7, 8);'
```

# Classifications

```bash
mysql --defaults-file="${HOME}"/replica.my.cnf -h tools-db s52585__cb -s -r -e 'select concat("INSERT INTO user_classification VALUES (null, -1, \"Import from report data\", ", if(status=7,1,0), ", ", new_id, ");") from reports inner join vandalism on id=revertid where status in (7, 8, 9);'
```
