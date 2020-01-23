#!/usr/bin/python3
# -*- coding: utf-8 -*-
import subprocess as shell
import os
from Crypto.PublicKey import RSA
from Crypto import Random
from Crypto.Cipher import PKCS1_OAEP

f = open('./pem/private.pem', 'r')
privateKey = RSA.importKey(f.read().encode('utf-8'))
f.close()
f = open('./pem/public.pem', 'r')
publicKey = RSA.importKey(f.read().encode('utf-8'))
f.close()

message = 'testmessage'
publicCipher = PKCS1_OAEP.new(publicKey)
cipherText = publicCipher.encrypt(message.encode('utf-8'))

print(cipherText)

privateCipher = PKCS1_OAEP.new(privateKey)
message2 = privateCipher.decrypt(cipherText).decode('utf-8')

print(message2)