# Custom Credential Integration Guide: Farmer ID (Testa Gava)

This guide uses the Verifiable ID (VerifiableId) as a base template to ensure the credential is automatically displayed by the Web Portal, overcoming its strict template loading mechanism.

## Prerequisites

The walt.id identity stack is running via docker-compose.

Droplet/external IP is 139.59.15.151.

Issuer DID is: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5Iiwia2lkIjoieW56SzZ1NTVTak82aEZFc1cwa0JLb25fYnB2cGY1enJyLVEzRk5IZUFWRSIsIngiOiJlM0NFMUVPcFl0RV82VXlJTjU4VUp3V21HR2VzVjNrWkhNVlpBQklRSTNNIn0 (Testa Gava).

Step 1: Define the Credential Templates (Web Portal & VC Repo)
To make the credential visible in the Web Portal (port 7102), we must create a custom file.

1.1 Create the Custom Template Structure (FarmerCredential.json)
Create this file to define the data fields for your new credential.

```json
{
"type": [
"VerifiableCredential",
"FarmerCredential"
],
"credentialSubject": {
"given_name": "",
"family_name": "",
"farm_name": "",
"issuing_authority": "Testa Gava",
"farm_type": "",
"license_no": "",
"region": "",
"status": ""
},
"issuer": {
"id": "",
"image": {
"id": "https://cdn-icons-png.flaticon.com/512/2010/2010178.png",
"type": "Image"
},
"name": "Testa Gava"
},
"expirationDate": ""
}
```

To ensure visibility in the Web Portal, we use the content of the custom file above, but save it with the name of a commonly loaded template (VerifiableId.json).

Action: Copy the content of FarmerCredential.json to the directory: waltid-applications/waltid-web-portal/.

Crucial Update: Modify the type array in VerifiableId.json to include "VerifiableId" so the Issuer API accepts it when using the generic ID: