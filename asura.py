from pixcat import Image
import os
import sys
import subprocess
import requests
from bs4 import BeautifulSoup as bs
import csv
import re
import io
import time
import shutil
import blessed

pathName = os.getcwd()
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:103.0) Gecko/20100101 Firefox/103.0'
}

url = 'https://asura.gg/'

# Image("0.jpg").fit_screen(enlarge=True).show()


def choose(options: list, limit: bool = True):
    oxot = []
    for option in options:
        oxot.append(f'"{option}"')
    options = ' '.join(oxot)
    if not limit:
        result = os.popen(f"gum choose --no-limit {options}")
    else:
        result = os.popen(f"gum choose {options}")
    if not limit:
        return result.read().replace(r"\n", "\n").strip().split("\n")
    else:
        return result.read().replace(r"\n", "\n").strip()


def get_html(url):
    r = requests.get(url, headers=headers)
    soup = bs(r.text, 'html.parser')
    return soup


def get_titles(soup):
    titles = []
    reader = soup.find_all('div', class_='utao styletwo')

    for uta in reader:
        alink = uta.find('a')
        # href = alink.get('href')
        src = alink.find('img')
        src = src.get('src')
        title = alink.get('title')
        titles.append(title)
    return titles


def get_images(soup):
    images = []
    reader = soup.find('div', class_='rdminimal')
    paragraphs = reader.find_all('p')
    for x in paragraphs:
        img = x.find('img')
        img = img.get('src').rstrip()
        img = img.strip()
        images.append(img)
        # print(img.get('src').rstrip())
    return images


def to_csv(res):
    with open("results.csv", "w") as f:
        write = csv.writer(f)
        write.writerow(res)


def to_file(image_links):
    with open("results.txt", "w", newline='\n') as f:
        f.write('\n'.join(image_links))
        # f.write(image_links)


def sanitize_link(title):

    # Remove all non-word characters (except numbers and letters)
    title = re.sub(r"[^\w\s]", '', title)
    # Replace all runs of whitespace with a single dash
    title = re.sub(r"\s+", '-', title)
    # ensure string is lower case
    title = title.lower()
    return title


def get_chapters(manga):
    title = ''
    # take urlname and create a valid url return single num of total chapters
    mangaUrl = f'https://www.asurascans.com/manga/{manga}/'
    print(mangaUrl)
    soup = get_html(mangaUrl)
    title = soup.find('span', class_='epcur epcurlast')
    title = title.text
    title = re.findall("\d+", title)[0]
    title = int(title)
    title = range(1, title + 1)
    print(title)
    return title


def download(images: list, manga_name: str, chapter):
    iterate = 0
    pathos = f'{pathName}/{manga_name}/'
    if not os.path.exists(f'{pathName}/{manga_name}/'):
        os.mkdir(f'{pathName}/{manga_name}/')
    elif not os.path.exists(f'{pathos}/{chapter}/'):
        os.mkdir(f'{pathos}/{chapter}/')
    # mkdir here
    paths = f'{pathos}/{chapter}'
    for img in images:
        iterate += 1
        res = requests.get(img, stream=True)
        name = f'img-{iterate}'
        if res.status_code == 200:
            with open(f'{paths}{name}', "wb") as f:
                shutil.copyfileobj(res.raw, f)
            print('Image downloaded successfully')
        else:
            print(f'Image {img} couldnt be downloaded')


def pixct(images):
    for i in images:
        Image(i).show()


def refresh(image):
    os.system("clear")
    Image(image).show()


term = blessed.Terminal()


def pixcat(images):
    with term.cbreak(), term.hidden_cursor():
        val = ''
        i = 0
        while val.lower() != 'q':
            val = term.inkey()
            if val.name == "KEY_DOWN":
                i += 1
                refresh(images[i])
            if val.name == "KEY_UP":
                i -= 1
                refresh(images[i])
            if val.name == "D":
                return


def pix(images):
    image_index = 0
    Image(images[image_index]).show()
    key = input()
    while True:

        if key == "j":
            image_index += 1
        elif key == "k":
            image_index += -1
        elif key == "q":
            sys.exit(0)

        if image_index == len(images):
            image_index = 0
        elif image_index < 0:
            image_index = len(images) - 1


def get_url(manga, ch_choice):
    mangaUrl = f'https://www.asurascans.com/{manga}-chapter-{ch_choice}/'
    return mangaUrl


def confirm(text: str):
    args = [f'gum confirm "{text}"']
    res = subprocess.Popen(args, shell=True)
    res.wait()
    rco = res.returncode
    if rco == 0:
        return True
    else:
        return False


def main():
    # get front page, parse titles, return chosen title.
    html = get_html(url)
    titles = get_titles(html)
    choice = choose(titles, limit=True)
    # turn title into a url friendly string
    urlname = sanitize_link(choice)
    # get max released chapters, return range. Then choose chapter.
    chapters = get_chapters(urlname)
    ch_choice = choose(chapters, limit=True)
    # stringify manhwa name and chapter for url
    manga_url = get_url(urlname, ch_choice)
    # get html for chapter, then get images and write to file.
    ch_html = get_html(manga_url)
    images = get_images(ch_html)
    to_file(images)
    pixcat(images)
    downl = confirm("Download Chapter?")
    if downl:
        download(images, urlname, ch_choice)
    else:
        main()
    # download(images, urlname, ch_choice)
    # pix(images)


main()
