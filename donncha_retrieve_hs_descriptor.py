"""
Retrieve hidden service descriptors from the
responsible HSDir's. This is a very rough
first copy and probably has bugs.

- Donncha O'Cearbhaill - donncha@donncha.is
  PGP: 0xAEC10762
"""

from time import time
from base64 import b32encode, b32decode
from hashlib import sha1
from struct import pack, unpack
from stem.descriptor import DocumentHandler
from stem.descriptor.remote import DescriptorDownloader
import argparse
from bisect import bisect_left
import urllib

# When provided with a Tor hidden service 'service_id', this script should output
# the predicted desc_id's which will be used to publish the HS descriptors for this
# HS into the future.

# Returns base_32 encode desc_id - descriptor-id = H(permanent-id | H(time-period | descriptor-cookie | replica))
def rend_compute_v2_desc_id(service_id_base32, replica, time = int(time()), descriptor_cookie = ""):#
  service_id = b32decode(service_id_base32, 1)
  time_period = get_time_period(time, 0, service_id)
  secret_id_part = get_secret_id_part_bytes(time_period, descriptor_cookie, replica)
  desc_id = rend_get_descriptor_id_bytes(service_id, secret_id_part)
  return b32encode(desc_id).lower()

# Calculates time period - time-period = (current-time + permanent-id-byte * 86400 / 256) / 86400
def get_time_period(time, deviation, service_id):
  REND_TIME_PERIOD_V2_DESC_VALIDITY = 24 * 60 * 60
  return int(((time + ((unpack('B', service_id[0])[0] * REND_TIME_PERIOD_V2_DESC_VALIDITY) ) / 256) ) / REND_TIME_PERIOD_V2_DESC_VALIDITY + deviation)

# Calculate secret_id_part - secret-id-part = H(time-period | descriptor-cookie | replica)
def get_secret_id_part_bytes(time_period, descriptor_cookie, replica):
  secret_id_part = sha1()
  secret_id_part.update(pack('>I', time_period)[:4]);
  if descriptor_cookie:
    secret_id_part.update(descriptor_cookie)
  secret_id_part.update('{0:02X}'.format(replica).decode('hex'))
  return secret_id_part.digest()

def rend_get_descriptor_id_bytes(service_id, secret_id_part):
  descriptor_id = sha1()
  descriptor_id.update(service_id)
  descriptor_id.update(secret_id_part)
  return descriptor_id.digest()

def find_responsible_HSDir(descriptor_id, consensus):
  fingerprint_list = []
  for _, router in consensus.routers.items():
    if "HSDir" in router.flags:
      fingerprint_list.append(router.fingerprint.decode("hex"))
  fingerprint_list.sort()

  descriptor_position = bisect_left(fingerprint_list, b32decode(descriptor_id, 1))

  responsible_HSDirs = []
  for i in range(0, 3):
    fingerprint = fingerprint_list[descriptor_position + i]
    router = consensus.routers[fingerprint.encode("hex").upper()]
    responsible_HSDirs.append({
      'nickname': router.nickname,
      'fingerprint': router.fingerprint,
      'address': router.address,
      'dir_port': router.dir_port,
      'descriptor_id': descriptor_id
    })
    
  return responsible_HSDirs

def main():
  REPLICAS = 2
  
  parser = argparse.ArgumentParser(description="This tool allows you to retrieve a copy of the raw hidden " \
                                               "service descriptor from the responsible hidden service " \
                                               "directories. It can also try all the responsible nodes to" \
                                               "determine how many will correctly return the descriptor.")
  parser.add_argument('onion_address', help='The hidden service address - e.g. (idnxcnkne4qt76tg.onion)')
  parser.add_argument("-v", "--verbose", action="store_true",
                  help="Show responsible HSDir's and try retrieve descriptors")
  args = parser.parse_args()
  
  responsible_HSDirs = []
  
  if args.verbose:
    print "Running in verbose mode"

  downloader = DescriptorDownloader()
  consensus = downloader.get_consensus(document_handler = DocumentHandler.DOCUMENT).run()[0]

  service_id, tld = args.onion_address.split(".")
  if tld == 'onion' and len(service_id) == 16 and service_id.isalnum():   
      for replica in range(0, REPLICAS):
        descriptor_id = rend_compute_v2_desc_id(service_id, replica, time())
        responsible_HSDirs.extend(find_responsible_HSDir(descriptor_id, consensus))
      
      # Loop through all the responsible HSDir's
      descriptor = ""
      for router in responsible_HSDirs:
        if (args.verbose == False) and descriptor:
          break

        if not router['dir_port']: continue

        url = 'http://'+router['address']+':'+str(router['dir_port'])+'/tor/rendezvous2/'+router['descriptor_id']

        if args.verbose:
          print url

        f = urllib.urlopen(url)
        if args.verbose:
          if f.getcode() == 200:
            descriptor = f.read().decode('utf-8')
          print b32decode(router['descriptor_id'], True).encode('hex')
          print str(f.getcode()) + '\t' + router['descriptor_id'] + '\t' + router['fingerprint'] + '\t' + router['nickname']
          
        else: # Loop until we find descriptor or error
          if f.getcode() == 200:
            descriptor = f.read().decode('utf-8')
            
      if descriptor:
        print descriptor
      else:
        print "[!] No descriptor could be retrieved for that onion address"
  else:
    print "[!] The onion address you provided is not valid"

if __name__ == '__main__':
    main()
