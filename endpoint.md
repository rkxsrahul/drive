# List of API Endpoints #

## 1. Healthz ##

```
Get /healthz
```

## 2. Create Pool ##

```
POST /v1/pool
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
name    string    required
```

## 3. List pools ##

```
GET /v1/pool
Header:-
Authorization: Bearer <<admin_token>>
```

## 4. Pool Details ##

```
GET /v1/pool/:pool_id
Header:-
Authorization: Bearer <<admin_token>>
```

## 5. Edit Pool ##

```
PUT /v1/pool/:pool_id
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
name    string    required
```

## 6. Delete Pool ##

```
DELETE /v1/pool/:pool_id
Header:-
Authorization: Bearer <<admin_token>>
```

## 7. Create Question in pool ##

```
POST /v1/question/:pool
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

UrlBody :-
pool -> pool_id

Body:-
{
	"title": "string, required",
	"type": "string, required",
	"options": [{
		"value":"string, required",
		"type":"string, required",
		"is_correct":boolean, required,
		"images_url":"string",
	}],
	"images_url":"string",
}

Note:- type value can only be "images" or "string"
```

## 8. List Questions in pool ##

```
GET /v1/question/:pool
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
pool -> pool_id
```

## 9. View Question Details ##

```
GET /v1/question/:pool/:id
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
pool -> pool_id
id -> question_id
```

## 10. Edit Question ##

```
PUT /v1/question/:pool/:id
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

UrlBody :-
pool -> pool_id
id -> question_id

Body:-
{
	"title": "string, required",
	"type": "string, required",
	"options": [{
		"value":"string, required",
		"type":"string, required",
		"is_correct":boolean, required,
		"images_url":"string",
	}],
	"images_url":"string",
}

Note:- type value can only be "images" or "string"
```

## 11. Delete Question ##

```
DELETE /v1/question/:pool/:id
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
pool -> pool_id
id -> question_id
```

## 12. Upload Image ##

```
POST /v1/upload_image
Content-Type: form-data
Header:-
Authorization: Bearer <<admin_token>>

Body:-
image file required
```

## 13. Delete Image ##

```
DELETE /v1/delete_image
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
url string required
```

## 14. Upload Questions with CSV ##

```
POST /v1/csv_question/:pool
Content-Type: form-data
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
pool -> pool_id
Body:-
questions file required

Note:- only csv file supported please convert excel sheet to csv then upload.
```

## 15. Create Test ##

```
POST /v1/test
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
{
	"name": "string, required",
	"duration": "string, required",
	"pools": [{
		"poolId":"string, required",
		"noOfQuestions":"integer, required",
	}]
}
```

## 16. List Tests ##

```
GET /v1/test
Header:-
Authorization: Bearer <<admin_token>>
```

## 17. Test Details ##

```
GET /v1/test/:test_id
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
test -> test_id
```

## 18. Edit Test ##

```
PUT /v1/test/:test_id
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

UrlBody :-
test -> test_id

Body:-
{
	"name": "string, required",
	"duration": "string, required",
	"pools": [{
		"poolId":"string, required",
		"noOfQuestions":"integer, required",
	}]
}
```

## 19. Delete Test ##

```
DELETE /v1/test/:test_id
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
test -> test_id
```

## 20. Create Job Team ##

```
POST /v1/teamjob
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
name    		string    required
description    	string    required
```

## 21. List job teams ##

```
GET /v1/teamjob
Header:-
Authorization: Bearer <<admin_token>>
```

## 22. Team Job Details ##

```
GET /v1/teamjob/:id
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
id -> team_id
```

## 23. Edit Job Team ##

```
PUT /v1/teamjob/:id
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

UrlBody :-
id -> team_id
Body:-
name    			string    required
description   string    required
```

## 24. Delete Job Team ##

```
DELETE /v1/teamjob/:id
Header:-
Authorization: Bearer <<admin_token>>
UrlBody :-
id -> team_id
```

## 25. Create Job ##

```
POST /v1/job
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
name    	string		required
summary 	string		required
location 	string		required
body 			string		required
teamId 		string		required
teamName 	string		required
skills 		[]string	required
```

## 26. List Jobs ##

```
GET /v1/job

Note:- all the keys initial will be in small letter
```

## 27. List Jobs by Teams ##

```
GET /v1/job/:team
UrlBody :-
team -> team_id

Note:- all the keys initial will be in small letter
```

## 28. Job Details ##

```
GET /v1/job/:team/:id

UrlBody :-
team -> team_id
id -> job_id

Note:- all the keys initial will be in small letter
```

## 29. Edit Job ##

```
PUT /v1/job/:team/:id
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

UrlBody :-
team -> team_id
id -> job_id
Body:-
name    	string		required
summary 	string		required
location 	string		required
body 			string		required
teamId 		string		required
teamName 	string		required
skills 		[]string	required
```

## 30. Delete Job ##

```
DELETE /v1/job/:team/:id
Header:-
Authorization: Bearer <<admin_token>>
UrlBody :-
team -> team_id
id -> job_id
```

## 31. Create College ##

```
POST /v1/college
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
name    		string    required
location    string    required
```

## 32. List Colleges on both side##

```
GET /v1/college

Query:-
page 		int		required
search	string
```

## 33. College Details ##

```
GET /v1/college/:college_id
Header:-
Authorization: Bearer <<admin_token>>
```

## 34. Edit College ##

```
PUT /v1/college/:college_id
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
name    		string    required
location    string    required
```

## 35. Delete College ##

```
DELETE /v1/college/:college_id
Header:-
Authorization: Bearer <<admin_token>>
```

## 36. Upload Colleges with CSV ##

```
POST /v1/csv_college
Content-Type: form-data
Header:-
Authorization: Bearer <<admin_token>>

Body:-
colleges file required

Note:- only csv file supported please convert excel sheet to csv then upload.
```

## 37. Create Drive ##

```
POST /v1/drive
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

Body:-
type				string	required
name 				string	required
start				string	required
end 				string	required
test_id			string	required
college_id 	string
```

## 38. List Drives ##

```
GET /v1/drive
Header:-
Authorization: Bearer <<admin_token>>

Params:-
value	string

Note:- if value is "ongoing" then you will fetch only ongoing drives otherwise all drives
```

## 39. Drive Details ##

```
GET /v1/drive/:drive_id
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
drive_id -> id of drive
```

## 40. Edit Drive ##

```
PUT /v1/drive/:drive_id
Header:-
Authorization: Bearer <<admin_token>>
Content-Type: application/json

UrlBody :-
drive_id -> id of drive

Body:-
type				string	required
name 				string	required
start				string	required
end 				string	required
test_id			string	required
college_id 	string
```

## 41. Delete Drive ##

```
DELETE /v1/drive/:drive_id
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
drive_id -> id of drive
```

## 42. Add user in Drive ##

```
POST /v1/drive/:drive_id/user
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
drive_id -> id of drive
Body:-
email string required
```

## 43. List users in a drive ##

```
GET /v1/drive/:drive_id/user
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
drive_id -> id of drive
```

## 44. Delete User from drive ##

```
DELETE /v1/drive/:drive_id/user/:email
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
drive_id -> id of drive
email -> emailId of user
```

## 45. Add users in  a drive through csv ##

```
POST /v1/drive/:drive_id/csv_user
Content-Type: form-data
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
drive_id -> id of drive
Body:-
users file required

Note:- only csv file supported please convert excel sheet to csv then upload.
```

## 46. List drives on user dashboard ##

```
GET /v1/user/drives
Header:-
Authorization: Bearer <<token>>
```

## 47. Generate new token for test and check user belongs to this test ##

```
GET /v1/user/drives/:drive_id/:testid
Header:-
Authorization: Bearer <<token>>

UrlBody :-
drive_id -> id of drive
testid -> id of test
```

## 48. Result api for fetching result of all users ##

```
GET /v1/result/drive/:id
Header:-
Authorization: Bearer <<admin_token>>

Params:-
id -> drive id
Query:-
page -> current page number
```

## 49. Pool wise result api of specific user ##

```
GET /v1/result/drive/:id
Header:-
Authorization: Bearer <<admin_token>>

Params:-
id -> drive id

Query Data:-
email -> email of user
```

## 50. List of completed tests of a user ##

```
GET /v1/user/completed
Header:-
Authorization: Bearer <<token>>
```

## 51. Test details at user id ##

```
GET /v1/user/test/:test_id
Header:-
Authorization: Bearer <<token>>

UrlBody :-
test -> test_id
```

## 52. Delete all users from a drive ##

```
DELETE /v1/drive/:id/user
Header:-
Authorization: Bearer <<admin_token>>

Params:-
id -> drive id
```

## 53. Delete user details from drive portal ##

```
DELETE /v1/user
Header:-
Authorization: Bearer <<admin_token>>

Query Data:-
email -> email of user
```

## 54. Fetch Questions in CSV ##

```
GET /v1/csv_question/:pool
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
pool -> pool_id
```

## 55. Fetch Result in CSV ##

```
GET /v1/result/drive/:id/csv
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
id -> drive_id
```

## 56. Summary of a drive ##

```
GET /v1/drive/:id/summary
Header:-
Authorization: Bearer <<admin_token>>

UrlBody :-
id -> drive_id
```

## 57. To download pool question's csv ##

```
GET /v1/pool_csv/:file

Params:-
file -> filename with extension
```

## 58. To download result's csv ##

```
GET /v1/result_csv/:file

Params:-
file -> filename with extension
```

## 59. All Drive result ##

```
GET /v1/result/drivewise
Header:-
Authorization: Bearer <<admin_token>>
```

## 60. Candidate result in a drive ##

```
GET /v1/result/drivewise/:drive
Header:-
Authorization: Bearer <<admin_token>>

Parameter :-
drive -> drive id
Query :-
page -> next page number (by default 1)
```

## 61. All Pool result ##

```
GET /v1/result/poolwise
Header:-
Authorization: Bearer <<admin_token>>
```

## 62. Candidates result in a pool ##

```
GET /v1/result/poolwise/:pool
Header:-
Authorization: Bearer <<admin_token>>

Parameter :-
pool -> pool id
Query :-
page -> next page number (by default 1)
```

## 63. Top 10 list in a drive ##

```
GET /v1/result/analytical/:drive
Header:-
Authorization: Bearer <<admin_token>>

Parameter :-
drive -> drive id
```

## 64. To disable or enable Mail Service##

```
PUT /v1/admin/mail/:value
Header:-
Authorization: Bearer <<admin token>>

Param data
value -> true if want to diable service and false if you want to enable service
```
