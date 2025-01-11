import csv
import os.path
import re
import sqlite3


def import_stuff():
    #steam_trade_url_regex = r"/^(https?:\/\/)?steamcommunity\.com\/tradeoffer\/new\/\?partner=[0-9]+&token=[a-zA-Z0-9_-]+$/igm"
    global youtube_url_regex
    youtube_url_regex = r"^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube(?:-nocookie)?\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|live\/|v\/)?)([\w\-]+)([^ ()]+)?$"

    good_rows = []

    bad_rows = []
    imported_pairs = []

    print('')
    print("Importing from 'input.csv'...")

    # Get input
    with open('input.csv', 'r', newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',', quotechar='|')
        count = -1

        # Validate youtube urls
        for row in reader:
            if count > 1:
                result = validate_row(row, count)

                if result == None:
                    good_rows.append(row)
                else:
                    bad_rows.append((row, result, count))

            count += 1

        con = sqlite3.connect("output.db")
        cur = con.cursor()

        for row in good_rows:            
            imported_pairs.append((extract_id_from_url(row[3]), row[1]))
        
        print(f'Import process completed')
        print(f'Total ready to enter db: {len(imported_pairs)}')
        print(f'Total which failed validation: {len(bad_rows)}')
        print('')

        return (imported_pairs, bad_rows)



def extract_id_from_url(url):
    split = url.split("/")
    id = split[-1].split('&')[0].split('?si=')[0].split('?t=')[0].split('?feature=shared')[0]
    if 'watch?v=' in url:
        id = id[8:]
    
    return id


def validate_row(row, count):
    row[-1] = row[-1].strip()
    url = row[-1]
    reason = ''
    is_valid_url = re.search(youtube_url_regex, url, re.IGNORECASE)
    
    if is_valid_url == None:
        return f"row {count}{' ' * (4 - len(str(count)))}: Bad url: '{url}'"
        
    if 'clip' in url or 'euri' in url:
        return f"row {count}{' ' * (4 - len(str(count)))}: Bad row: Weird/clip url. Manually inspect"

    if len(row) > 4:
        return f"row {count}{' ' * (4 - len(str(count)))}: Bad row: Good url, but bad row size. Manually inspect."

    return None 


def db_setup():
    con = sqlite3.connect("output.db")
    cur = con.cursor()
    cur.execute("""CREATE TABLE IF NOT EXISTS videos (id INTEGER PRIMARY KEY, url TEXT UNIQUE, uploader_username TEXT);""")
    cur.execute("""CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, ip TEXT UNIQUE);""")
    cur.execute("""CREATE TABLE IF NOT EXISTS votes (user_id INTEGER NOT NULL, video_url TEXT NOT NULL, score INTEGER NOT NULL);""")
    cur.execute("""CREATE TABLE IF NOT EXISTS active_votes (
		user_id INTEGER PRIMARY KEY NOT NULL,
		start_time INTEGER NOT NULL,
		a TEXT,
		b TEXT
	);""")



def export_stuff(imported_pairs):
    db_setup()

    print(f"\nExporting {len(imported_pairs)} submissions to 'output.db' ")

    con = sqlite3.connect("output.db")
    cur = con.cursor()

    #for (id, name) in imported_pairs: 
        #pass
        #print(f'{id}, {name}')      
    cur.executemany(f"INSERT OR IGNORE INTO videos (url, uploader_username) VALUES (?, ?)", imported_pairs)
        #cur.execute(f"INSERT OR IGNORE INTO videos (url, uploader_username) VALUES ('{id}', '{name.replace("'", "''")}')")
    
    con.commit()
    print('Export completed\n')



def list_good_rows(imported_pairs):
    if len(imported_pairs) > 0:
        print('')

    count = 2
    for (id, name) in imported_pairs:
        print(f"row {count}{' ' * (4 - len(str(count)))}: {id}{' ' * (33 - len(id))} |  {name}")
        count += 1
    
    print('')
    print(f'Total: {len(imported_pairs)}')
    print('')


def list_bad_rows(bad_rows):
    print('')

    for (row, reason, count) in bad_rows:
        print(reason)
    
    print('')
    print(f'Total: {len(bad_rows)}')
    print('')


def resolver(bad_rows):
    new_bad_rows = []
    db_setup()

    con = sqlite3.connect("output.db")
    cur = con.cursor()

    print('\nStarting resolver...\n')
    print('The resolver will iterate through all failed submissions and request valid information.')
    print('This fixed data will directly exported to the output database.\n')
    print("Enter 'skip' to skip a submission, 'exit' to exit the loop.")
    print('')


    for (row, reason, count) in bad_rows:

        while True:
            print('')
            print(reason)
            print('Raw submission data:')
            print(f'     {' | '.join(row)}')
            print('')
            print('Manually input new valid in the format of "<username> <new url>"')

            action = input(f"CSP Resolver | Row {count} >> ")

            if action == 'skip':
                new_bad_rows.append((row, reason))
                print('skipped')
                break
            
            if action == 'exit':
                print('exiting...\n')
                return bad_rows
            
            split = action.split(' ')

            if not len(split) == 2:
                print('Invalid input. Must be "<username> <new url>"')
                continue

            row = ['', split[0].strip(), '', split[1].strip()]
            result = validate_row(row, 'fail')

            if result == None:
                cur.execute(f"INSERT OR IGNORE INTO videos (url, uploader_username) VALUES (?, ?)", (extract_id_from_url(row[3]), row[1]))
                #cur.execute(f"INSERT OR IGNORE INTO videos (url, uploader_username) VALUES ('{extract_id_from_url(split[1].strip())}', '{split[0].strip().replace("'", "''")}')")
                con.commit()
                print("Valid data accepted. Successfully added to 'output.db'")
                break
            else:
                print(result)
    
    return new_bad_rows



def db_to_txt():
    con = sqlite3.connect("output.db")
    cur = con.cursor()

    res = cur.execute(f"SELECT url FROM videos")

    with open('id_dump.txt', 'w') as file:
        for row in res:
            file.write(row[0] + '\n')

    print('\ndumped\n')




def main():
    imported_pairs = []
    bad_rows = []

    print('')
    print('Clip Submission Processor by Mythitorium')
    print('')
    print('Commands:')
    print('  import')
    print('  export')
    print('  listgood')
    print('  listbad')
    print('  resolver : Iterate over and resolve submissions which failed auto-validation')
    print('  db2txt : Takes everything in the database and dumps them into a text file')
    print('  exit')
    print('')
    while True:
        action = input("CSP >> ")

        match action:
            case 'import': (imported_pairs, bad_rows) = import_stuff()
            case 'export': export_stuff(imported_pairs)
            case 'listgood': list_good_rows(imported_pairs)
            case 'listbad': list_bad_rows(bad_rows)
            case 'resolver' : bad_rows = resolver(bad_rows)
            case 'db2txt' : db_to_txt()
            case 'exit': break
            case _: print('\nCommand is le invalid\n')


if __name__=="__main__":
    main()
