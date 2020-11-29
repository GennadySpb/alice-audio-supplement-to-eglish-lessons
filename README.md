What is it?!
===
Educational project to have fun, bring profit (children use it when doing homework), also get familiar with serverless and create a skill for Alice.


Audio Source
===
Original audio files from [site](https://xn----dtbhthpdbkkaet.xn--p1ai/El_uchebniki/Elektronnie_resursi/multimediynye-prilozheniya-k-uchebnikam/angliyskiy-yazyk-/index.php).
The files are an audio attachment to the second grade English textbook.


Extract audio resource ID
===
After upload and processing of audio files you should get UUID of them
```bash
curl  -H "Authorization: OAuth $DIALOG_TOKEN" \
   "https://dialogs.yandex.net/api/v1/skills/${SKILL_ID}/sounds" > sounds.json
cat sounds.json | jq '[ .sounds | sort_by(.originalName) | .[] | { (.originalName[20:22]) : .id} ]' | fgrep \"
``` 

Data Preparation
===
Marking data partially made via ffmpeg and service and [STT service](https://cloud.yandex.ru/docs/speechkit/stt/request).

Something like
```bash
# convert first 7 sec to OGG format
for f in *.mp3; 
do 
    ffmpeg -i $f -c:a libvorbis -q:a 4 -to 7  "${f/%mp3/ogg}"; 
done

# output file ID and result of TTS
for f in *.ogg; 
do 
    foo=${f##Brilliant_RS_1_1502_}; 
    echo ${foo%%.ogg} -- $(curl -X POST -s \
     -H "Authorization: Bearer ${IAM_TOKEN}" \
     -H "Transfer-Encoding: chunked"  \
     --data-binary "@${f}" \
     "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize?topic=general&folderId=${FOLDER_ID}") ; 
done
```
