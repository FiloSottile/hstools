"""
List a Hidden Service HSDir for a given consensus or at the current instant.

- Donncha O'Cearbhaill - donncha@donncha.is
- Filippo Valsorda
"""

from time import mktime, time
from base64 import b32encode, b32decode
from hashlib import sha1
from struct import pack, unpack
from stem.descriptor import parse_file, DocumentHandler
from stem.descriptor.remote import DescriptorDownloader
import argparse
from bisect import bisect_left

# Returns base_32 encode desc_id - descriptor-id = H(permanent-id | H(time-period | descriptor-cookie | replica))
def rend_compute_v2_desc_id(service_id_base32, replica, time, descriptor_cookie = ""):#
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
  
  parser = argparse.ArgumentParser()
  parser.add_argument('onion_address', help='The hidden service address - e.g. (idnxcnkne4qt76tg.onion)')
  parser.add_argument('--consensus', help='The optional consensus file', required=False)
  args = parser.parse_args()

  if args.consensus is None:
    downloader = DescriptorDownloader()
    consensus = downloader.get_consensus(document_handler = DocumentHandler.DOCUMENT).run()[0]
    t = time()
  else:
    with open(args.consensus) as f:
      consensus = next(parse_file(f, 'network-status-consensus-3 1.0', document_handler = DocumentHandler.DOCUMENT))
    t = mktime(consensus.valid_after.timetuple())

  service_id, tld = args.onion_address.split(".")
  if tld == 'onion' and len(service_id) == 16 and service_id.isalnum():   
      for replica in range(0, REPLICAS):
        descriptor_id = rend_compute_v2_desc_id(service_id, replica, t)
        print descriptor_id + '\t' + b32decode(descriptor_id, True).encode('hex')
        for router in find_responsible_HSDir(descriptor_id, consensus):
          print router['fingerprint'] + '\t' + router['nickname']

  else:
    print "[!] The onion address you provided is not valid"

if __name__ == '__main__':
    main()
