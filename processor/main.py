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
    id = split[-1]
    if not 'clip' in url:
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
            print('Manually input new valid in the format of "<username>|<new url>"')

            action = input(f"CSP Resolver | Row {count} >> ")

            if action == 'skip':
                new_bad_rows.append((row, reason))
                print('skipped')
                break
            
            if action == 'exit':
                print('exiting...\n')
                return bad_rows
            
            split = action.split('|')

            if not len(split) == 2:
                print('\nInvalid input. Must be "<username>|<new url>"')
                print("-----")
                continue

            new_row = ['', split[0].strip(), '', split[1].strip()]
            #result = validate_row(new_row, 'fail')

            cur.execute(f"INSERT OR IGNORE INTO videos (url, uploader_username) VALUES (?, ?)", (extract_id_from_url(new_row[3]), new_row[1]))
            #cur.execute(f"INSERT OR IGNORE INTO videos (url, uploader_username) VALUES ('{extract_id_from_url(split[1].strip())}', '{split[0].strip().replace("'", "''")}')")
            con.commit()
            print("\nValid data accepted. Successfully added to 'output.db'")
            print("-----")
            break
    
    return new_bad_rows



def db_to_txt():
    con = sqlite3.connect("output.db")
    cur = con.cursor()

    res = cur.execute(f"SELECT url FROM videos ORDER BY id")

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




def voting_distrib(exclude_culled: bool):
    pass

    clips = {}

    con = sqlite3.connect("live.db")
    cur = con.cursor()

    # Db query 
    res = []
    if exclude_culled:
        res = cur.execute("""
            SELECT video_url, score FROM votes WHERE 
                          video_url NOT IN (SELECT url FROM culled_videos)
                          AND video_url IN (SELECT url FROM videos)
        """)
    else:
        res = cur.execute('SELECT video_url, score FROM votes WHERE video_url IN (SELECT url FROM videos)')
    
    # process rows
    for row in res:
        url = row[0]
        score = row[1]
        if not url in clips:
            clips[url] = { "w": 0, "l" : 0 }
        else:
            if score == 1:
                clips[url]["w"] += 1
            else:
                clips[url]["l"] += 1

    # Set ratios  
    for key in clips:
        clips[key]["ratio"] = ((clips[key]["w"]/(clips[key]["w"]+clips[key]["l"])) * 100)

    # Sort by ratio
    clips = dict(sorted(clips.items(), key=lambda x:x[1]['ratio']))

    # Print results
    count = 0
    additive_total = 0
    vote_total = 0
    for key in clips:
        data = clips[key]
        wins = data["w"]
        loss = data["l"]
        ratio = data['ratio']
        additive_total += (wins/(wins+loss)) * 100
        vote_total += wins + loss
        #print(f'{count} | {key} -- ratio: {(ratio)}% - {wins+loss} votes - {wins} wins & {loss} losses) ')
        #print(f'{969 - count} | ratio: {ratio}% -- link: <https://youtu.be/{key}>')
        res = cur.execute(f"""SELECT uploader_username FROM videos WHERE url = '{key}'""")
        user = ''
        for row in res:
            user = row[0]
        print(f'new Clip("#{952 - count} {user}", "{key}", {round(ratio, 3)}),')
        count += 1
    
    # Print totals
    if exclude_culled:
        print('POOL: UNCULLED CLIPS')
    else:
        print('POOL: UNCULLED + CULLED CLIPS')

    print(f'AVERAGE WINRATE: {round(additive_total/count)-13}%')
    print(f'TOTAL VOTES: {vote_total}')
    print(f'AVERAGE VOTES PER CLIP: {round(vote_total/count)-13}')
    print(f'total clips in this pool: {count-13}')
    
    # Add another number if you want to see data at a different thresh
    for thresh in [45, 50, 55, 60, 65, 70, 75, 80, 90]:
        total = 0
        for key in clips:
            if clips[key]['ratio'] >= thresh:
                total += 1
        
        print(f'Amount of clips above {thresh}%: {max(0, total-13)}')
    
    print('\nNOTE: TOTALS ARE ALL OFFSET BY -13')
    

if __name__=="__main__":
    #voting_distrib(False)
#
#
    #items = [
    #    0,
    #    83,
    #    84,
    #    16,
    #    99,
    #    15,
    #    100,
    #    8,
    #    23,
    #    92,
    #    76,
    #    91,
    #    7,
    #    3,
    #    80,
    #    87,
    #    19,
    #    96,
    #    12,
    #    11,
    #    20,
    #    95,
    #    78,
    #    88,
    #    4,
    #    1,
    #    82,
    #    85,
    #    17,
    #    98,
    #    14,
    #    101,
    #    9,
    #    22,
    #    93,
    #    77,
    #    90,
    #    6,
    #    2,
    #    81,
    #    86,
    #    18,
    #    97,
    #    13,
    #    10,
    #    21,
    #    94,
    #    79,
    #    89,
    #    5,
    #]
    #items.sort()

    #print(items)
    main()