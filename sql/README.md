# Legacy data import from report interface

# Edits marked as in progress from report

```bash
mysql --defaults-file="${HOME}"/replica.my.cnf -h tools-db s52585__cb -s -r -e\
'SELECT CONCAT("INSERT INTO edit VALUES (", new_id, ", 0, 3) ON DUPLICATE KEY UPDATE id=id; '\
'INSERT INTO edit_edit_group VALUES (", new_id, ", 2);") '\
'FROM reports INNER JOIN vandalism ON id=revertid WHERE status IN (7, 8);' > data.edit-set.2.sql
````

# Edits from report status

```bash
mysql --defaults-file="${HOME}"/replica.my.cnf -h tools-db s52585__cb -s -r -e\
'SELECT CONCAT("INSERT INTO edit VALUES (", new_id, ", 0, ", if(status=7,1,0), ") ON DUPLICATE KEY UPDATE id=id; '\
'INSERT INTO edit_edit_group VALUES (", new_id, ", 1); '\
'INSERT INTO user_classification VALUES (null, -1, \"Import from report data\", ", if(status=7,1,0), ", ", new_id, ") ON DUPLICATE KEY UPDATE id=id;") '\
'FROM reports INNER JOIN vandalism ON id=revertid WHERE status IN (7, 8);' > data.edit-set.1.sql
````

# Edits from original training sets

```bash
function parse_and_import_xml {
  dataset_id=$1; file_path=$2

  echo "Processing ${dataset_id} -- ${file_path}";
  grep -E '(<EditID>[0-9]+</EditID>)|(<isVandalism>(true|false)</isVandalism>)' "${file_path}" | \
  sed -e 's/\s*<EditID>//' \
      -e 's/<\/EditID>//' \
      -e 's/\s*<isVandalism>false<\/isVandalism>/1/' \
      -e 's/\s*<isVandalism>true<\/isVandalism>/0/' | \
  while read edit_id;
  do
    read -r status_id
    echo "INSERT INTO edit VALUES (${edit_id}, 0, ${status_id}) ON DUPLICATE KEY UPDATE id=id; " \
         "INSERT INTO edit_edit_group VALUES (${edit_id}, ${dataset_id}); " \
         "INSERT INTO user_classification VALUES (null, -1, \"Import from original data\", ${status_id}, ${edit_id}) ON DUPLICATE KEY UPDATE id=id; "
  done > "sql/data.edit-set.${dataset_id}.sql"
}

parse_and_import_xml "7" "cluebotng/editsets/C/train.xml"
parse_and_import_xml "8" "cluebotng/editsets/C/trial.xml"

parse_and_import_xml "9" "cluebotng/editsets/D/train.xml"
parse_and_import_xml "10" "cluebotng/editsets/D/trial.xml"
parse_and_import_xml "11" "cluebotng/editsets/D/bayestrain.xml"
parse_and_import_xml "12" "cluebotng/editsets/D/all.xml"

parse_and_import_xml "13" "cluebotng-testing/editsets/C/train.xml"
parse_and_import_xml "14" "cluebotng-testing/editsets/C/trial.xml"

parse_and_import_xml "15" "cluebotng-testing/editsets/D/train.xml"
parse_and_import_xml "16" "cluebotng-testing/editsets/D/trial.xml"
parse_and_import_xml "17" "cluebotng-testing/editsets/D/bayestrain.xml"
parse_and_import_xml "18" "cluebotng-testing/editsets/D/all.xml"

parse_and_import_xml "19" "cluebotng-testing/editsets/Auto/train.xml"
parse_and_import_xml "20" "cluebotng-testing/editsets/Auto/trial.xml"

parse_and_import_xml "21" "cluebotng-testing/editsets/OldTriplet/train.xml"
parse_and_import_xml "22" "cluebotng-testing/editsets/OldTriplet/trial.xml"
parse_and_import_xml "23" "cluebotng-testing/editsets/OldTriplet/bayestrain.xml"
parse_and_import_xml "24" "cluebotng-testing/editsets/OldTriplet/all.xml"

parse_and_import_xml "25" "cluebotng-testing/editsets/RandomEdits50-50/train.xml"
parse_and_import_xml "26" "cluebotng-testing/editsets/RandomEdits50-50/trial.xml"
parse_and_import_xml "27" "cluebotng-testing/editsets/RandomEdits50-50/all.xml"

parse_and_import_xml "28" "cluebotng-testing/editsets/VeryLarge/train.xml"
parse_and_import_xml "29" "cluebotng-testing/editsets/VeryLarge/trial.xml"
parse_and_import_xml "30" "cluebotng-testing/editsets/VeryLarge/bayestrain.xml"
parse_and_import_xml "31" "cluebotng-testing/editsets/VeryLarge/all.xml"
```
