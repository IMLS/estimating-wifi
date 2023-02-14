import random
import requests
from fpdf import FPDF
import pandas as pd
import admin_tool_kit_checks as ck

EXPECTED_HEADERS = ['fscs_id', 'label', 'address']

def genKeyword():
    """
    Create a unique key for each sensor in word-word-word format
    """
    lines = open("5000-more-common.txt").read() 
    line = lines[0:] 
    words = line.split()
    mypass = ""
    for x in range(0, 4):
        myword = random.choice(words).upper()
        if x == 3:
            mypass += myword
        else:
            mypass += myword + "-" 
    #print passphrase
    return mypass
    
def testHB(s, sensor, key):
    """
    Call "beat_the_heart" to test new sensor
    """
    token_url = s + "/rpc/login"
    body = {"fscs_id": sensor, "api_key": key}
    tr = requests.post(token_url, json=body)
    t0 = tr.json()["token"]
    url = s + "/rpc/beat_the_heart" 
    body = {
        "_sensor_serial": "abcde",
        "_sensor_version": "01.0",
    }
    headers = {"Authorization": f"Bearer {t0}"}
    r = requests.post(url, json=body, headers=headers)
    print("BEAT RESPONSE", r.json())

def genPDF(sensor, key, label):
        """
        Create and save a PDF
        """
        instructions='pdfs/instructions.txt'

        #Setup PDF object
        pdf = FPDF()  
        pdf.add_page()
        pdf.set_font("Arial", size = 15)
        #Title cell
        pdf.cell(200, 10, txt = "SETUP INSTRUCTIONS FOR LIBRARY",
                ln = 1, align = 'C')
        pdf.cell(200, 10, txt = "SENSOR: " + sensor,
                ln = 2, align = 'C')
        pdf.cell(200, 10, txt = "PASS: " + key,
                ln = 3, align = 'C')
        pdf.cell(200, 10, txt = "LABEL: " + label,
                ln = 4, align = 'C')
        pdf.cell(200, 10, txt = '',
                ln = 4, align = 'C')
        # pdf.multicell(200, 10, txt=data, align='C')

        with open(instructions, 'r+') as f:
            for line in f: 
                pdf.multi_cell(200, 5, txt=line, align='L')
        
        # save the pdf with name .pdf
        pdf.output("pdfs/" + sensor + ".pdf")  
        print("PDF created for ", sensor)

# This is essentially the entire CSV checker.
# It is pulled out into the util file so that it can be unit tested.
# All of the code in `libadmin` should be pulled out in a similar way, so that
# the functions can all be unit tested.
# Instead, many of them have *parts* that are tested, but not the top-level
# commands themselves. Check might serve as a model for how that could be done.
# Essentially, all calls to `sys.exit()` need to be removed, and replaced with
# `return` statements. This makes testing possible.
# Some judicious try/except statements might be needed as well.
def check(filename):
    does_file_exist = ck.check_file_exists(filename)
    if not does_file_exist:
        print("File '{}' does not exist.".format(filename))
        return -1
    does_filename_end_with = ck.check_filename_ends_with(filename, "csv")
    if not does_filename_end_with:
        print("{} does not end with CSV.".format(filename))
        return -1
    # Read in the CSV with headers
    df = pd.read_csv(filename, header=0)
    # check_headers will throw specific errors for specific mismatches.
    r3 = ck.check_any_nulls(df)
    if len(r3) != 0:
        for r in r3:
            print("We're missing data in column '{}'".format(r))
        return -1
    # https://stackoverflow.com/questions/30487993/python-how-to-check-if-two-lists-are-not-empty
    # Checking lists involves truthiness and falsiness of []. I'll keep it simple.
    # And, more importantly... make sure it works. I'll check the list length.
    # FIXME: This should be a *set comparison*, which would solve all 
    # of these length checks. Set difference should yield the empty set.
    r1 = ck.check_headers(df, EXPECTED_HEADERS)
    if isinstance(r1, int):
        return r1
    elif isinstance(r1, list) and (len(r1) != 0):
        for r in r1:
            print("Expected header '{}', found '{}'.".format(r["expected"], r["actual"]))
        return -1
    r2 = ck.check_library_ids(df)
    if len(r2) != 0:
        for r in r2:
            print("{} is not a valid library ID.".format(r))
        return -1
    
    return 0
