import re
from pathlib import Path
# Pulled from check.py



# https://stackoverflow.com/questions/82831/how-do-i-check-whether-a-file-exists-without-exceptions
# Checks to see if a file exists in the local filesystem.
def check_file_exists(filename):
    file = Path(filename)
    return file.is_file()

# https://stackoverflow.com/questions/10873777/in-python-how-can-i-check-if-a-filename-ends-in-html-or-files
# Checks the filename for a particular extension.
def check_filename_ends_with(filename, ext):
    return filename.endswith(ext)

# Checks that the headers on a dataframe are what we expect.
# This should be rewritten as a set intersection problem: 
# take the set of headers in the DF, the set that are expected, subtract,
# and if the set is non-empty, we have a problem.
def check_headers(df, expected_headers):
    results = []
    actual_headers = list(df.columns.values)
    if len(actual_headers) > len(expected_headers):
        print("CSV has more headers than expected.")
        return -1
    if len(actual_headers) < len(expected_headers):
        print("CSV has fewer headers than expected.")
        return -1
    for expected, actual in zip(expected_headers, actual_headers):
        if not (expected == actual):
            results.append({"expected": expected, "actual": actual})
    return results

# https://www.dataquest.io/wp-content/uploads/2019/03/python-regular-expressions-cheat-sheet.pdf
# https://www.pythoncheatsheet.org/cheatsheet/regular-expressions
# Checks to see that all the FSCS ids in the dataframe are of the correct form.
# FIXME: This checks for AA0000, not AA0000-001. The latter pass through, but the check
# is against a simpler ID.
def check_library_ids(df):
    results = []
    regex = re.compile('[A-Z]{2}[0-9]{4}') 
    ids = list(df['fscs_id'].values)
    for id in ids:
        if not regex.match(id):
            print("{} did not match".format(id))
            results.append(id)
    return results

# Checks for any nulls in the dataframe.
# Should be extended to also look for empty strings.
def check_any_nulls(df):
    found_nulls = []
    actual_headers = list(df.columns.values) 
    for header in actual_headers:
        column = list(df.isna()[header])
        # https://stackoverflow.com/questions/35784074/does-python-have-andmap-ormap
        anymap = any(column)
        if anymap:
            found_nulls.append(header)
    return found_nulls