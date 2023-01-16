import csv
import random
import sys, getopt
import requests
from fpdf import FPDF

class setupSensor:

    def __init__(self, s) -> None:
        self.sensor = s
        self.key = ""
        self.label = ""
    
    def callSetup(self, s, u, p):
        """
        Call "sensor_setup" to save info to the database.
        """
        token_url = s + "/rpc/login"
        body = {"fscs_id": u, "api_key": p}
        # Need to post, not get, if you're passing params in the body.
        tr = requests.post(token_url, json=body)
        t0 = tr.json()["token"]
        url = s + "/rpc/sensor_setup" 
        body = {
            "_sensor": self.sensor,
            "_key": self.key,
            "_label": self.label,
        }
        headers = {"Authorization": f"Bearer {t0}"}
        r = requests.post(url, json=body, headers=headers)
        print("INSERT RESPONSE", r.json())
    
    def setLabel(self, l):
        """
        Set labels for the sensor object
        """
        self.label = l 

    def genKeyword(self):
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
        self.key = mypass

    def testHB(self, s):
        """
        Call "beat_the_heart" to test new sensor
        """
        token_url = s + "/rpc/login"
        body = {"fscs_id": self.sensor, "api_key": self.key}
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

    def genPDF(self):
        """
        Create and save a PDF
        """
        #Setup PDF object
        pdf = FPDF()  
        pdf.add_page()
        pdf.set_font("Arial", size = 15)
        #Title cell
        pdf.cell(200, 10, txt = "SETUP INSTRUCTIONS FOR LIBRARY",
                ln = 1, align = 'C')
        pdf.cell(200, 10, txt = "SENSOR: " + self.sensor,
                ln = 2, align = 'C')
        pdf.cell(200, 10, txt = "PASS: " + self.key,
                ln = 3, align = 'C')
        pdf.cell(200, 10, txt = "LABEL: " + self.label,
                ln = 4, align = 'C')
        
        # save the pdf with name .pdf
        pdf.output(self.sensor + ".pdf")  
        print("PDF created for ", self.sensor)


def readCSV(filename, server, u, p):
    """
    Loop to read the CSV and setup each sensor
    """
    with open(filename, newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            #print(row)
            #Create the sensor object
            print("Setting up Sensor ", row['fscs_id'])
            s = setupSensor(row['fscs_id'])
            print("Inserting Sensor into DB")
            #Create unique passkey
            s.genKeyword()
            #Set sensor label
            s.setLabel(row['label'])
            #Call the database setup function
            s.callSetup(server, u, p)
            #Test beat_the_heart
            print("Testing Sesnor ", row['fscs_id'])
            s.testHB(server)
            #Create PDF 
            print("Building PDF for ", row['fscs_id'])
            s.genPDF()
            #Complete message and itterate loop
            print("***SENSOR COMPLETE***")
            

def main(argv):
    """
    Main function to handle commands line args
    """
    filename = ''
    server_addr = ''
    user = ''
    password = ''
    opts, args = getopt.getopt(argv,"hf:s:u:p:",["file=","server=","user=", "pass="])
    for opt, arg in opts:
        if opt == '-h':
            print ('test_admin_setup.py -f <csv_path> -s <protocol://url:port> -u <username> -p <password>')
            sys.exit()
        elif opt in ("-f", "--file"):
            filename = arg
        elif opt in ("-s", "--server"):
            server_addr = arg
        elif opt in ("-u", "--user"):
            user = arg
        elif opt in ("-p", "--pass"):
            password = arg
    # print ('Args f ', filename)
    # print ('Args s ', server_addr)
    # print ('Args u ', user)
    # print ('Args p ', password)

    print("***READING CSV FILE*** ", filename)
    readCSV(filename, server_addr, user, password)
    print("***SETUP COMPLETE***")

if __name__ == "__main__":
    main(sys.argv[1:])