package force

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

SObject//SObjectResponse received from force.com API after insert of an sobject.
type SObjectResponse struct {
	ID      string    `force:"id,omitempty"`
	Errors  ApiErrors `force:"error,omitempty"` //TODO: Not sure if ApiErrors is the right object
	Success bool      `force:"success,omitempty"`
}

//DescribeSObjects describes all SObjects
func (forceAPI *ForceApi) DescribeSObjects() (map[string]*SObjectMetaData, error) {
	if err := forceAPI.getApiSObjects(); err != nil {
		return nil, err
	}

	return forceAPI.apiSObjects, nil
}

//DescribeSObject describes a object by given name
func (forceAPI *ForceApi) DescribeSObject(objectName string) (resp *SObjectDescription, err error) {
	// Check cache
	resp, ok := forceAPI.apiSObjectDescriptions[objectName]
	if !ok {
		// Attempt retrieval from api
		sObjectMetaData, ok := forceAPI.apiSObjects[objectName]
		if !ok {
			err = fmt.Errorf("Unable to find metadata for object: %v", objectName)
			return
		}

		uri := sObjectMetaData.URLs[sObjectDescribeKey]

		resp = &SObjectDescription{}
		err = forceAPI.Get(uri, nil, resp)
		if err != nil {
			return
		}

		// Create Comma Separated String of All Field Names.
		// Used for SELECT * Queries.
		length := len(resp.Fields)
		if length > 0 {
			var allFields bytes.Buffer
			for index, field := range resp.Fields {
				// Field type location cannot be directly retrieved from SQL Query.
				if field.Type != "location" {
					if index > 0 && index < length {
						allFields.WriteString(", ")
					}
					allFields.WriteString(field.Name)
				}
			}

			resp.AllFields = allFields.String()
		}

		forceAPI.apiSObjectDescriptions[objectName] = resp
	}

	return
}

//GetSObject returns object by given id, name and wanted fields
func (forceAPI *ForceApi) GetSObject(id string, objectName string, fields []string, out interface{}) error {
	uri := strings.Replace(forceAPI.apiSObjects[objectName].URLs[rowTemplateKey], idKey, id, 1)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	return forceAPI.Get(uri, params, out)
}

//InsertSObject creates new object by given name and object
func (forceAPI *ForceApi) InsertSObject(objectName string, in interface{}) (resp *SObjectResponse, err error) {
	uri := forceAPI.apiSObjects[objectName].URLs[sObjectKey]

	resp = &SObjectResponse{}
	err = forceAPI.Post(uri, nil, in, resp)

	return
}

//UpdateSObject updates object by given id,name and object
func (forceAPI *ForceApi) UpdateSObject(id string, objectName string, in interface{}) error {
	uri := strings.Replace(forceAPI.apiSObjects[objectName].URLs[rowTemplateKey], idKey, id, 1)

	return forceAPI.Patch(uri, nil, in, nil)
}

//DeleteSObject deletes object by given id and name
func (forceAPI *ForceApi) DeleteSObject(id string, objectName string) error {
	uri := strings.Replace(forceAPI.apiSObjects[objectName].URLs[rowTemplateKey], idKey, id, 1)

	return forceAPI.Delete(uri, nil)
}

//GetSObjectByExternalID Not sure if or how this works. Please contact me if something explodes.
func (forceAPI *ForceApi) GetSObjectByExternalID(id string, externalID string, objectName string, fields []string, out interface{}) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[objectName].URLs[sObjectKey],
		externalID, id)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	err = forceAPI.Get(uri, params, out)

	return
}

//UpsertSObjectByExternalID Not sure if or how this works. Please contact me if something explodes.
func (forceAPI *ForceApi) UpsertSObjectByExternalID(id string, externalID string, objectName string, in interface{}) (resp *SObjectResponse, err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[objectName].URLs[sObjectKey],
		externalID, id)

	resp = &SObjectResponse{}
	err = forceAPI.Patch(uri, nil, in, resp)

	return
}

//DeleteSObjectByExternalID Not sure if or how this works. Please contact me if something explodes.
func (forceAPI *ForceApi) DeleteSObjectByExternalID(id string, externalID string, objectName string) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[objectName].URLs[sObjectKey],
		externalID, id)

	err = forceAPI.Delete(uri, nil)

	return
}
