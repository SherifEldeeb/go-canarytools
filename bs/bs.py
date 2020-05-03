import requests
from bs4 import BeautifulSoup
import re
import json

page = requests.get("https://f1a7c45a.canary.tools/apidocs/incident_attributes.html")
soup = BeautifulSoup(page.content, 'html.parser')
allpre = soup.find_all('pre')
allpre.pop(0)
for i in allpre:
    t = i.get_text()
    lines = [re.sub(r"\/\*[^\*]*\*\/","", l) for l in t.splitlines()]
    e = {}
    e["description"] = lines[0].split("=")[1].replace("\"", "").strip().title()
    e["object_name"] = e["description"].replace(" ", "")+"Event"
    e["logtype"] = lines[1].split("=")[1].replace("\"", "").strip().title()
    e["dict"] = "".join([l.strip().strip("{") for l in lines[2:-1]]).split("=")[1].strip().replace(",", "\n").splitlines()
    e["dict"] = [re.sub(r'^"([^"]*).*',r"\1 string", l).title() for l in e["dict"]]
    e["dict_text"] = "\n".join(e["dict"])

    the_text = '''package canarytokens
    // {object_name} is a {description} event, they are part of incidents
    type {object_name} struct {{
        {dict_text}
    }}
    '''.format(**e)

    print(the_text)

