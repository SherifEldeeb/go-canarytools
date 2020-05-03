import requests
from bs4 import BeautifulSoup
import re
import json

tmp = '''
package canarytokens

// {object_name}Event is a {description} event 
type {object_name}Event struct {{
    {dict_text}
}}
'''
whitelist = set('abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ')
line = "{key} {otype} // {comment}"

def getType(t):
    lines = [l.strip() for l in t.splitlines() if l.strip() != ""]
    tmpDict = {}
    for i in lines:
        k = i.split(" ")[0]
        v = " ".join(i.split(" ")[1:])
        tmpDict[k] = v
    return tmpDict

page = requests.get("https://canarytools.readthedocs.io/en/latest/incident_attributes.html")
soup = BeautifulSoup(page.content, 'html.parser')
allattr = soup.find_all('ul', {"class":"first last simple"})
for ul in allattr:
    d = getType(ul.get_text())
    with open("d:\\demofile2.txt", "a") as f:
        for i in d.items():
            f.write("{0} {1}\n".format(*i))
            print("{0} {1}".format(*i))


    # e = {}
    # e["dict_text"] = ""
    # for li in ul.findChildren("li"):
    #     line = "{k} {otype} // {comment}\n"
    #     tmpDict = {}

    #     t = li.get_text()
    #     try:
    #         k = li.findChildren("strong")[0].get_text()
    #     except:
    #         pass

    #     v = t[len(k)+3:].strip("\"")
    #     tmpDict["otype"] = "string"
    #     if "(dict)" in k:
    #         tmpDict["otype"] = "[]interface{{}}"
    #     k = k.replace("(str)", "").replace("(dict)", "").strip()
    #     if k == "description":
    #         e[k.title()] = ''.join(filter(whitelist.__contains__, v)).title()
    #     if k == "type":
    #         object_name = v.replace("-", " ").title().replace(" ", "")
    #         e["object_name"] = object_name
    # print(e)
