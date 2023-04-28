import click 
import admin_tool_kit as tk
import admin_tool_libraries as tl
import admin_tool_sensors as ts
import pandas as pd
import os
from dotenv import load_dotenv 

load_dotenv()

SERVER_ADDR=os.getenv('SERVER_ADDR')
ADMIN_USERNAME=os.getenv('ADMIN_USERNAME')
ADMIN_PASSWORD=os.getenv('ADMIN_PASSWORD')


@click.group()
def cli():
    pass

@cli.command(name='validateCSV')
@click.option('--file', required=True, help='File path to csv.')
def validateCSV(file):
    """A program to validate the CSV before use."""
    if tk.check(file) == 0:
        print("valid CSV")

@cli.command(name='readLog')
def readLog():
    """A program to read the logfile. NOT FUNCTIONAL"""
    print("reading log")

@cli.command(name='addSingleLibrary')
@click.option('--library', prompt=True, help='fscs_id for new library')
def addSingleLibrary(library):
    """A program to add a single library to the database."""
    tl.callSetup(SERVER_ADDR, ADMIN_USERNAME, ADMIN_PASSWORD, library)

@cli.command(name='addMultipleLibraries')
@click.option('--username', prompt=True, help='Admin account username')
@click.password_option('--password', help='Admin account password')
@click.option('--file', required=True, help='File path to csv with libraries.')
def addMultipleLibraries():
    """A program to add multiple libraries from a CSV file. NOT FUNCTIONAL"""
    print("adding mutliple libraries")

@cli.command(name='updateLibrary')
@click.option('--library', prompt=True, help='fscs_id for library to be updated')
@click.option('--address', help='New address for existing sensor')
def updateLibrary():
    """A program to update an existing library. NOT FUNCTIONAL"""
    print("updating library")

@cli.command(name='removeLibrary')
@click.option('--library', prompt=True, help='fscs_id for library to be removed')
def removeLibrary(library):
    """A program to remove an existing library"""
    tl.callRemoveLibrary(SERVER_ADDR, ADMIN_USERNAME, ADMIN_PASSWORD, library)

@cli.command(name='addSingleSensor')
# @click.option('--username', prompt=True, help='Admin account username')
# @click.password_option('--password', help='Admin account password')
@click.option('--sensor', prompt=True, help='fscs_id for new sesnor')
@click.option('--label', prompt=True, help='label for new sesnor')
def addSingleSesnor(sensor, label):
    """A program to add a single sensor to the database."""
    tempKey = tk.genKeyword()
    ts.callSetup(SERVER_ADDR, ADMIN_USERNAME, ADMIN_PASSWORD, sensor, tempKey, label)
    tk.testHB (SERVER_ADDR, sensor, tempKey)
    tk.genPDF(sensor, tempKey, label)

@cli.command(name='addMultipleSensors')
@click.option('--username', prompt=True, help='Admin account username')
@click.password_option('--password', help='Admin account password')
@click.option('--file', required=True, help='File path to csv.')
def addMultipleSensors(username, password, file):
    """A program to add multiple sensors from a CSV file. NOT FUNCTIONAL."""
    if tk.check(file) == 0:
        df = pd.read_csv(file, header=0)
        print(df)
    else:
        print("error adding mutliple sessnors")

@cli.command(name='updateSensor')
@click.option('--username', prompt=True, help='Admin account username')
@click.password_option('--password', help='Admin account password')
@click.option('--sensor', prompt=True, help='fscs_id for sensor to be updated')
@click.option('--password', help='New password for existing sensor')
@click.option('--labels', help='New label for existing sensor')
@click.option('--address', help='New address for existing sensor')
@click.option('--genPDF', is_flag=True, show_default=True, default=False, help='Generate a new PDF for existing sensor')
def updateSensor():
    """A program to update an existing sensor. NOT FUNCTIONAL"""
    print("updating password")

@cli.command(name='removeSensor')
# @click.option('--username', prompt=True, help='Admin account username')
# @click.password_option('--password', help='Admin account password')
@click.option('--sensor', prompt=True, help='fscs_id for new sesnor to be removed')
def removeSensor(sensor):
    """A program to remove an existing sensor"""
    ts.callRemoveSensor(SERVER_ADDR, ADMIN_USERNAME, ADMIN_PASSWORD, sensor)

@cli.command(name='checkDB')
@click.option('--username', prompt=True, help='Admin account username')
@click.password_option('--password', help='Admin account password')
@click.option('--server', prompt=True, help='The IP/Hostname of the PiSpots server')
def removeSensor():
    """A program to test the PiSpots stack. NOT FUNCTIONAL"""
    print("testing stack")

if __name__ == '__main__':
    cli()