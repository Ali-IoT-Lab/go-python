#!/usr/bin/python3
 
import io
import sys
sys.stdout = io.TextIOWrapper(sys.stdout.buffer,encoding='utf-8')

print(sys.getfilesystemencoding())

s = "中文"  #py3 默认就是utf8 unicode向下兼容utf-8
#s_to_gpk = s.encode("gbk")
#print(s_to_gpk.decode("gbk")) #py3是不需要转换的，但是要转的话还要转回去
print("语言学习")
