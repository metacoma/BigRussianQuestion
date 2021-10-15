from __future__ import print_function
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support import ui
import sys
import time
import re
import os

def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)

#print("STARTED")

driver = webdriver.Firefox()
driver.get(os.environ['OTVET_URL'])
body = driver.find_element_by_tag_name('body')
wait = ui.WebDriverWait(driver ,0.5)

n = 0

time.sleep(5)
sended = {} 

ignore = {
  "Прочее компьютерное": 1,
  "Прочие": 1,
  "Железо": 1,
  "Другие языки и технологии": 1,
  "Клиентские": 1,
  "Программное обеспечение": 1,
  "Прочее компьютерное": 1,
  "Мобильная связь": 1,
  "Интернет": 1,
  "Прочее фото-видео": 1,
  "Кино, Театр": 1,
  "Музыка": 1,
  "Консольные": 1,
  "Прочие Авто-темы": 1,
  "Верстка, CSS, HTML, SVG": 1,
  "Естественные науки": 1,
  "Бухгалтерия, Аудит, Налоги": 1,
  "Прочие Авто-темы": 1,
  "Веб-дизайн": 1,
  "Офисная техника": 1,
  "Прочие юридические вопросы": 1,
  "Мобильные": 1,
  "Десерты, Сладости, Выпечка": 1,
  "Вторые блюда": 1,
  "Сервис, Обслуживание, Тюнинг": 1,
  "Мобильные устройства": 1,
  "Строительство и Ремонт": 1,
  "PHP": 1,
  "Python": 1,
  "Java": 1,
  "JavaScript": 1,
  "jQuery": 1,
  "SQL": 1,
  "Perl": 1,
  "C#": 1,
  "Ответы Mail.ru": 1,
  "C/C++": 1,
  "Почта Mail.ru": 1,
  "Геометрия": 1,
  "Android": 1,
  "iOS": 1,
  "Системное администрирование": 1,
  "ВУЗы, Колледжи": 1,
  "Другие предметы": 1,
  "Прочие юридические вопросы": 1,
  "История": 1,
  "Школы": 1,
  "Физика": 1,
  "Информатика": 1,
  "География": 1,
  "Алгебра": 1,
  "Техника для дома": 1,
  "Русский язык": 1,
  "Биология": 1,
  "Гражданское право": 1,
  "Иностранные языки": 1,
  "Математика": 1
  
} 


accept = driver.find_elements_by_xpath('//div[text() = " Продолжить просмотр "]')
if len(accept): 
  accept[0].click()

def need_skip(data):
  STOP_LIST = [".*Помогите найти.*", ".*(вопрос|фото) внутри.*", "79807422416", '.*см[\.]?\ внутри[\.]?', ".*см. вн.", '.*см вопрос', ".*\+СМ$", ".*см\+$"]
  for i in STOP_LIST:
    r = re.compile(i, re.IGNORECASE)
    if (r.match(data)):
      return 1
  return 0

while 1:

    show_more = driver.find_elements_by_xpath("//div[@size = 'promo']")

    if (len(show_more) > 0): 
      try:
        new_count = driver.find_elements_by_xpath("//div[@size='promo']/span/i")[0].get_attribute('innerText')
      except:   
        continue 
      #print(f'new count {new_count}')
      show_more[0].click() 
      time.sleep(3)
      show = 1

      whos = driver.find_elements_by_xpath("//div/div/div/div/div/div/div/div/span[1]/span/a") 
      categories = driver.find_elements_by_xpath("//div/div/div/div/div/div/div/div/span[2]/span/a") 
      #answers = driver.find_elements_by_xpath('//div/div/div/div/div/div/div/a[contains(@href, "/question/")]')
      answers = driver.find_elements_by_xpath('//div/div/div/div/div/div/div[2]/a[contains(@href, "/question/")]')
      #print('answers {0}, whos: {1}, categories: {2}'.format(len(answers), len(whos), len(categories)))
      for answer in answers:

          #print("FOUND") 
          who = whos[show - 1].get_attribute('innerText')
          category = categories[show - 1].get_attribute('innerText')


          question = answer.get_attribute('innerText')
          #print(f"FOUND question {question}") 

          # send_data = question + "(" + who + ", " + category + ")" 
          # send_data = question + "(" + category + ")" 
          send_data = question 
        
          if send_data not in sended:
              if category not in ignore and not need_skip(question) and len(question) >= 35:
                print(send_data,  flush=True)
                eprint(send_data,  flush=True)
                sended[send_data] = 1
              else:
                eprint("SKIP " + send_data,  flush=True)

          if (show == int(new_count)): 
              break
          show = show + 1
      body.send_keys(Keys.HOME)
            
    time.sleep(3)
        




#for n in range(0, 60):
#    time.sleep(5)
#    body.send_keys(Keys.END)
#
#    tweets = driver.find_elements_by_xpath("//div[@lang = 'ru']/span")
#    for tweet in tweets:
#        try:
#            print(tweet.get_attribute('innerText'))
#        except:
#            pass
