// The rds package provides types and functions for interaction with the AWS
// Relational Database service (rds)
package rds

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mitchellh/goamz/aws"
)

// The Rds type encapsulates operations operations with the Rds endpoint.
type Rds struct {
	aws.Auth
	aws.Region
	httpClient *http.Client
}

const APIVersion = "2013-09-09"

// New creates a new Rds instance.
func New(auth aws.Auth, region aws.Region) *Rds {
	return NewWithClient(auth, region, aws.RetryingClient)
}

func NewWithClient(auth aws.Auth, region aws.Region, httpClient *http.Client) *Rds {
	return &Rds{auth, region, httpClient}
}

func (rds *Rds) query(params map[string]string, resp interface{}) error {
	params["Version"] = APIVersion
	params["Timestamp"] = time.Now().In(time.UTC).Format(time.RFC3339)

	endpoint, err := url.Parse(rds.Region.RdsEndpoint)
	if err != nil {
		return err
	}

	sign(rds.Auth, "GET", "/", params, endpoint.Host)
	endpoint.RawQuery = multimap(params).Encode()
	r, err := rds.httpClient.Get(endpoint.String())

	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode > 200 {
		return buildError(r)
	}

	decoder := xml.NewDecoder(r.Body)
	decodedBody := decoder.Decode(resp)

	return decodedBody
}

func buildError(r *http.Response) error {
	var (
		err    Error
		errors xmlErrors
	)
	xml.NewDecoder(r.Body).Decode(&errors)
	if len(errors.Errors) > 0 {
		err = errors.Errors[0]
	}
	err.StatusCode = r.StatusCode
	if err.Message == "" {
		err.Message = r.Status
	}
	return &err
}

func multimap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}

func makeParams(action string) map[string]string {
	params := make(map[string]string)
	params["Action"] = action
	return params
}

// ----------------------------------------------------------------------------
// Rds objects

type DBInstance struct {
	Address                    string        `xml:"Endpoint>Address"`
	AllocatedStorage           int           `xml:"AllocatedStorage"`
	AvailabilityZone           string        `xml:"AvailabilityZone"`
	BackupRetentionPeriod      int           `xml:"BackupRetentionPeriod"`
	DBInstanceClass            string        `xml:"DBInstanceClass"`
	DBInstanceIdentifier       string        `xml:"DBInstanceIdentifier"`
	DBInstanceStatus           string        `xml:"DBInstanceStatus"`
	DBName                     string        `xml:"DBName"`
	Engine                     string        `xml:"Engine"`
	EngineVersion              string        `xml:"EngineVersion"`
	MasterUsername             string        `xml:"MasterUsername"`
	MultiAZ                    bool          `xml:"MultiAZ"`
	Port                       int           `xml:"Endpoint>Port"`
	PreferredBackupWindow      string        `xml:"PreferredBackupWindow"`
	PreferredMaintenanceWindow string        `xml:"PreferredMaintenanceWindow"`
	VpcSecurityGroupIds        []string      `xml:"VpcSecurityGroups>VpcSecurityGroupMembership>VpcSecurityGroupId"`
	DBSecurityGroupNames       []string      `xml:"DBSecurityGroups>DBSecurityGroup>DBSecurityGroupName"`
	DBSubnetGroup              DBSubnetGroup `xml:"DBSubnetGroup"`
	DBParameterGroupName       string        `xml:"DBParameterGroups>DBParameterGroup>DBParameterGroupName"`
}

type DBSecurityGroup struct {
	Description              string   `xml:"DBSecurityGroupDescription"`
	Name                     string   `xml:"DBSecurityGroupName"`
	EC2SecurityGroupNames    []string `xml:"EC2SecurityGroups>EC2SecurityGroup>EC2SecurityGroupName"`
	EC2SecurityGroupIds      []string `xml:"EC2SecurityGroups>EC2SecurityGroup>EC2SecurityGroupId"`
	EC2SecurityGroupOwnerIds []string `xml:"EC2SecurityGroups>EC2SecurityGroup>EC2SecurityGroupOwnerId"`
	EC2SecurityGroupStatuses []string `xml:"EC2SecurityGroups>EC2SecurityGroup>Status"`
	CidrIps                  []string `xml:"IPRanges>IPRange>CIDRIP"`
	CidrStatuses             []string `xml:"IPRanges>IPRange>Status"`
}

type DBSubnetGroup struct {
	Description string   `xml:"DBSubnetGroupDescription"`
	Name        string   `xml:"DBSubnetGroupName"`
	Status      string   `xml:"SubnetGroupStatus"`
	SubnetIds   []string `xml:"Subnets>Subnet>SubnetIdentifier"`
	VpcId       string   `xml:"VpcId"`
}

type DBSnapshot struct {
	AllocatedStorage     int    `xml:"AllocatedStorage"`
	AvailabilityZone     string `xml:"AvailabilityZone"`
	DBInstanceIdentifier string `xml:"DBInstanceIdentifier"`
	DBSnapshotIdentifier string `xml:"DBSnapshotIdentifier"`
	Engine               string `xml:"Engine"`
	EngineVersion        string `xml:"EngineVersion"`
	InstanceCreateTime   string `xml:"InstanceCreateTime"`
	Iops                 int    `xml:"Iops"`
	LicenseModel         string `xml:"LicenseModel"`
	MasterUsername       string `xml:"MasterUsername"`
	OptionGroupName      string `xml:"OptionGroupName"`
	PercentProgress      int    `xml:"PercentProgress"`
	Port                 int    `xml:"Port"`
	SnapshotCreateTime   string `xml:"SnapshotCreateTime"`
	SnapshotType         string `xml:"SnapshotType"`
	SourceRegion         string `xml:"SourceRegion"`
	Status               string `xml:"Status"`
	VpcId                string `xml:"VpcId"`
}

type DBParameterGroup struct {
	DBParameterGroupFamily string `xml:"DBParameterGroupFamily"`
	DBParameterGroupName   string `xml:"DBParameterGroupName"`
	Description            string `xml:"Description"`
}

type Parameter struct {
	ApplyMethod    string `xml:"ApplyMethod"`
	ParameterName  string `xml:"ParameterName"`
	ParameterValue string `xml:"ParameterValue"`
}

// ----------------------------------------------------------------------------
// Create

// The CreateDBInstance request parameters
type CreateDBInstance struct {
	AllocatedStorage           int
	AvailabilityZone           string
	BackupRetentionPeriod      int
	DBInstanceClass            string
	DBInstanceIdentifier       string
	DBName                     string
	DBSubnetGroupName          string
	Engine                     string
	EngineVersion              string
	Iops                       int
	MasterUsername             string
	MasterUserPassword         string
	MultiAZ                    bool
	Port                       int
	PreferredBackupWindow      string // hh24:mi-hh24:mi
	PreferredMaintenanceWindow string // ddd:hh24:mi-ddd:hh24:mi
	PubliclyAccessible         bool
	VpcSecurityGroupIds        []string
	DBSecurityGroupNames       []string
	DBParameterGroupName       string

	SetAllocatedStorage      bool
	SetBackupRetentionPeriod bool
	SetIops                  bool
	SetPort                  bool
}

func (rds *Rds) CreateDBInstance(options *CreateDBInstance) (resp *SimpleResp, err error) {
	params := makeParams("CreateDBInstance")

	if options.SetAllocatedStorage {
		params["AllocatedStorage"] = strconv.Itoa(options.AllocatedStorage)
	}

	if options.SetBackupRetentionPeriod {
		params["BackupRetentionPeriod"] = strconv.Itoa(options.BackupRetentionPeriod)
	}

	if options.SetIops {
		params["Iops"] = strconv.Itoa(options.Iops)
	}

	if options.SetPort {
		params["Port"] = strconv.Itoa(options.Port)
	}

	if options.AvailabilityZone != "" {
		params["AvailabilityZone"] = options.AvailabilityZone
	}

	if options.DBInstanceClass != "" {
		params["DBInstanceClass"] = options.DBInstanceClass
	}

	if options.DBInstanceIdentifier != "" {
		params["DBInstanceIdentifier"] = options.DBInstanceIdentifier
	}

	if options.DBName != "" {
		params["DBName"] = options.DBName
	}

	if options.DBSubnetGroupName != "" {
		params["DBSubnetGroupName"] = options.DBSubnetGroupName
	}

	if options.Engine != "" {
		params["Engine"] = options.Engine
	}

	if options.EngineVersion != "" {
		params["EngineVersion"] = options.EngineVersion
	}

	if options.MasterUsername != "" {
		params["MasterUsername"] = options.MasterUsername
	}

	if options.MasterUserPassword != "" {
		params["MasterUserPassword"] = options.MasterUserPassword
	}

	if options.MultiAZ {
		params["MultiAZ"] = "true"
	}

	if options.PreferredBackupWindow != "" {
		params["PreferredBackupWindow"] = options.PreferredBackupWindow
	}

	if options.PreferredMaintenanceWindow != "" {
		params["PreferredMaintenanceWindow"] = options.PreferredMaintenanceWindow
	}

	if options.PubliclyAccessible {
		params["PubliclyAccessible"] = "true"
	}

	for j, group := range options.VpcSecurityGroupIds {
		params["VpcSecurityGroupIds.member."+strconv.Itoa(j+1)] = group
	}

	for j, group := range options.DBSecurityGroupNames {
		params["DBSecurityGroups.member."+strconv.Itoa(j+1)] = group
	}

	if options.DBParameterGroupName != "" {
		params["DBParameterGroupName"] = options.DBParameterGroupName
	}

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// The CreateDBSecurityGroup request parameters
type CreateDBSecurityGroup struct {
	DBSecurityGroupName        string
	DBSecurityGroupDescription string
}

func (rds *Rds) CreateDBSecurityGroup(options *CreateDBSecurityGroup) (resp *SimpleResp, err error) {
	params := makeParams("CreateDBSecurityGroup")

	params["DBSecurityGroupName"] = options.DBSecurityGroupName
	params["DBSecurityGroupDescription"] = options.DBSecurityGroupDescription

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// The CreateDBSubnetGroup request parameters
type CreateDBSubnetGroup struct {
	DBSubnetGroupName        string
	DBSubnetGroupDescription string
	SubnetIds                []string
}

func (rds *Rds) CreateDBSubnetGroup(options *CreateDBSubnetGroup) (resp *SimpleResp, err error) {
	params := makeParams("CreateDBSubnetGroup")

	params["DBSubnetGroupName"] = options.DBSubnetGroupName
	params["DBSubnetGroupDescription"] = options.DBSubnetGroupDescription

	for j, group := range options.SubnetIds {
		params["SubnetIds.member."+strconv.Itoa(j+1)] = group
	}

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// The CreateDBSecurityGroup request parameters
type AuthorizeDBSecurityGroupIngress struct {
	Cidr                    string
	DBSecurityGroupName     string
	EC2SecurityGroupId      string
	EC2SecurityGroupName    string
	EC2SecurityGroupOwnerId string
}

func (rds *Rds) AuthorizeDBSecurityGroupIngress(options *AuthorizeDBSecurityGroupIngress) (resp *SimpleResp, err error) {
	params := makeParams("AuthorizeDBSecurityGroupIngress")

	if attr := options.Cidr; attr != "" {
		params["CIDRIP"] = attr
	}

	if attr := options.EC2SecurityGroupId; attr != "" {
		params["EC2SecurityGroupId"] = attr
	}

	if attr := options.EC2SecurityGroupOwnerId; attr != "" {
		params["EC2SecurityGroupOwnerId"] = attr
	}

	if attr := options.EC2SecurityGroupName; attr != "" {
		params["EC2SecurityGroupName"] = attr
	}

	params["DBSecurityGroupName"] = options.DBSecurityGroupName

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// The CreateDBParameterGroup request parameters
type CreateDBParameterGroup struct {
	DBParameterGroupFamily string
	DBParameterGroupName   string
	Description            string
}

func (rds *Rds) CreateDBParameterGroup(options *CreateDBParameterGroup) (resp *SimpleResp, err error) {
	params := makeParams("CreateDBParameterGroup")

	params["DBParameterGroupFamily"] = options.DBParameterGroupFamily
	params["DBParameterGroupName"] = options.DBParameterGroupName
	params["Description"] = options.Description

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// Describe

// DescribeDBInstances request params
type DescribeDBInstances struct {
	DBInstanceIdentifier string
}

type DescribeDBInstancesResp struct {
	RequestId   string       `xml:"ResponseMetadata>RequestId"`
	DBInstances []DBInstance `xml:"DescribeDBInstancesResult>DBInstances>DBInstance"`
}

func (rds *Rds) DescribeDBInstances(options *DescribeDBInstances) (resp *DescribeDBInstancesResp, err error) {
	params := makeParams("DescribeDBInstances")

	params["DBInstanceIdentifier"] = options.DBInstanceIdentifier

	resp = &DescribeDBInstancesResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DescribeDBSecurityGroups request params
type DescribeDBSecurityGroups struct {
	DBSecurityGroupName string
}

type DescribeDBSecurityGroupsResp struct {
	RequestId        string            `xml:"ResponseMetadata>RequestId"`
	DBSecurityGroups []DBSecurityGroup `xml:"DescribeDBSecurityGroupsResult>DBSecurityGroups>DBSecurityGroup"`
}

func (rds *Rds) DescribeDBSecurityGroups(options *DescribeDBSecurityGroups) (resp *DescribeDBSecurityGroupsResp, err error) {
	params := makeParams("DescribeDBSecurityGroups")

	params["DBSecurityGroupName"] = options.DBSecurityGroupName

	resp = &DescribeDBSecurityGroupsResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DescribeDBSubnetGroups request params
type DescribeDBSubnetGroups struct {
	DBSubnetGroupName string
}

type DescribeDBSubnetGroupsResp struct {
	RequestId      string          `xml:"ResponseMetadata>RequestId"`
	DBSubnetGroups []DBSubnetGroup `xml:"DescribeDBSubnetGroupsResult>DBSubnetGroups>DBSubnetGroup"`
}

func (rds *Rds) DescribeDBSubnetGroups(options *DescribeDBSubnetGroups) (resp *DescribeDBSubnetGroupsResp, err error) {
	params := makeParams("DescribeDBSubnetGroups")

	params["DBSubnetGroupName"] = options.DBSubnetGroupName

	resp = &DescribeDBSubnetGroupsResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DescribeDBSnapshots request params
type DescribeDBSnapshots struct {
	DBInstanceIdentifier string
	DBSnapshotIdentifier string
	SnapshotType         string
}

type DescribeDBSnapshotsResp struct {
	RequestId   string       `xml:"ResponseMetadata>RequestId"`
	DBSnapshots []DBSnapshot `xml:"DescribeDBSnapshotsResult>DBSnapshots>DBSnapshot"`
}

func (rds *Rds) DescribeDBSnapshots(options *DescribeDBSnapshots) (resp *DescribeDBSnapshotsResp, err error) {
	params := makeParams("DescribeDBSnapshots")

	if options.DBInstanceIdentifier != "" {
		params["DBInstanceIdentifier"] = options.DBInstanceIdentifier
	}

	if options.DBSnapshotIdentifier != "" {
		params["DBSnapshotIdentifier"] = options.DBSnapshotIdentifier
	}

	if options.SnapshotType != "" {
		params["SnapshotType"] = options.SnapshotType
	}

	resp = &DescribeDBSnapshotsResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DescribeDBParameterGroups request params
type DescribeDBParameterGroups struct {
	DBParameterGroupName string
}

type DescribeDBParameterGroupsResp struct {
	RequestId         string             `xml:"ResponseMetadata>RequestId"`
	DBParameterGroups []DBParameterGroup `xml:"DescribeDBParameterGroupsResult>DBParameterGroups>DBParameterGroup"`
}

func (rds *Rds) DescribeDBParameterGroups(options *DescribeDBParameterGroups) (resp *DescribeDBParameterGroupsResp, err error) {
	params := makeParams("DescribeDBParameterGroups")

	params["DBParameterGroupName"] = options.DBParameterGroupName

	resp = &DescribeDBParameterGroupsResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DescribeDBParameters request params
type DescribeDBParameters struct {
	DBParameterGroupName string
	Source               string
}

type DescribeDBParametersResp struct {
	RequestId  string      `xml:"ResponseMetadata>RequestId"`
	Parameters []Parameter `xml:"DescribeDBParametersResult>Parameters>Parameter"`
}

func (rds *Rds) DescribeDBParameters(options *DescribeDBParameters) (resp *DescribeDBParametersResp, err error) {
	params := makeParams("DescribeDBParameters")

	params["DBParameterGroupName"] = options.DBParameterGroupName

	if attr := options.Source; attr != "" {
		params["Source"] = attr
	}

	resp = &DescribeDBParametersResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DeleteDBInstance request params
type DeleteDBInstance struct {
	FinalDBSnapshotIdentifier string
	DBInstanceIdentifier      string
	SkipFinalSnapshot         bool
}

func (rds *Rds) DeleteDBInstance(options *DeleteDBInstance) (resp *SimpleResp, err error) {
	params := makeParams("DeleteDBInstance")

	params["DBInstanceIdentifier"] = options.DBInstanceIdentifier

	// If we don't skip the final snapshot, we need to specify a final
	// snapshot identifier
	if options.SkipFinalSnapshot {
		params["SkipFinalSnapshot"] = "true"
	} else {
		params["FinalDBSnapshotIdentifier"] = options.FinalDBSnapshotIdentifier
	}

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DeleteDBSecurityGroup request params
type DeleteDBSecurityGroup struct {
	DBSecurityGroupName string
}

func (rds *Rds) DeleteDBSecurityGroup(options *DeleteDBSecurityGroup) (resp *SimpleResp, err error) {
	params := makeParams("DeleteDBSecurityGroup")

	params["DBSecurityGroupName"] = options.DBSecurityGroupName

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DeleteDBSubnetGroup request params
type DeleteDBSubnetGroup struct {
	DBSubnetGroupName string
}

func (rds *Rds) DeleteDBSubnetGroup(options *DeleteDBSubnetGroup) (resp *SimpleResp, err error) {
	params := makeParams("DeleteDBSubnetGroup")

	params["DBSubnetGroupName"] = options.DBSubnetGroupName

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// DeleteDBParameterGroup request params
type DeleteDBParameterGroup struct {
	DBParameterGroupName string
}

func (rds *Rds) DeleteDBParameterGroup(options *DeleteDBParameterGroup) (resp *SimpleResp, err error) {
	params := makeParams("DeleteDBParameterGroup")

	params["DBParameterGroupName"] = options.DBParameterGroupName

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

type RestoreDBInstanceFromDBSnapshot struct {
	DBInstanceIdentifier    string
	DBSnapshotIdentifier    string
	AutoMinorVersionUpgrade bool
	AvailabilityZone        string
	DBInstanceClass         string
	DBName                  string
	DBSubnetGroupName       string
	Engine                  string
	Iops                    int
	LicenseModel            string
	MultiAZ                 bool
	OptionGroupName         string
	Port                    int
	PubliclyAccessible      bool

	SetIops bool
	SetPort bool
}

func (rds *Rds) RestoreDBInstanceFromDBSnapshot(options *RestoreDBInstanceFromDBSnapshot) (resp *SimpleResp, err error) {
	params := makeParams("RestoreDBInstanceFromDBSnapshot")

	params["DBInstanceIdentifier"] = options.DBInstanceIdentifier
	params["DBSnapshotIdentifier"] = options.DBSnapshotIdentifier

	if options.AutoMinorVersionUpgrade {
		params["AutoMinorVersionUpgrade"] = "true"
	}

	if options.AvailabilityZone != "" {
		params["AvailabilityZone"] = options.AvailabilityZone
	}

	if options.DBInstanceClass != "" {
		params["DBInstanceClass"] = options.DBInstanceClass
	}

	if options.DBName != "" {
		params["DBName"] = options.DBName
	}

	if options.DBSubnetGroupName != "" {
		params["DBSubnetGroupName"] = options.DBSubnetGroupName
	}

	if options.Engine != "" {
		params["Engine"] = options.Engine
	}

	if options.SetIops {
		params["Iops"] = strconv.Itoa(options.Iops)
	}

	if options.LicenseModel != "" {
		params["LicenseModel"] = options.LicenseModel
	}

	if options.MultiAZ {
		params["MultiAZ"] = "true"
	}

	if options.OptionGroupName != "" {
		params["OptionGroupName"] = options.OptionGroupName
	}

	if options.SetPort {
		params["Port"] = strconv.Itoa(options.Port)
	}

	if options.PubliclyAccessible {
		params["PubliclyAccessible"] = "true"
	}

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// ModifyDBParameterGroup request parameters
type ModifyDBParameterGroup struct {
	DBParameterGroupName string
	Parameters           []Parameter
}

func (rds *Rds) ModifyDBParameterGroup(options *ModifyDBParameterGroup) (resp *SimpleResp, err error) {
	params := makeParams("ModifyDBParameterGroup")

	params["DBParameterGroupName"] = options.DBParameterGroupName

	for j, group := range options.Parameters {
		params["Parameters.member."+strconv.Itoa(j+1)+".ApplyMethod"] = group.ApplyMethod
		params["Parameters.member."+strconv.Itoa(j+1)+".ParameterName"] = group.ParameterName
		params["Parameters.member."+strconv.Itoa(j+1)+".ParameterValue"] = group.ParameterValue
	}

	resp = &SimpleResp{}

	err = rds.query(params, resp)

	if err != nil {
		resp = nil
	}

	return
}

// Responses

type SimpleResp struct {
	RequestId string `xml:"ResponseMetadata>RequestId"`
}

type xmlErrors struct {
	Errors []Error `xml:"Error"`
}

// Error encapsulates an Rds error.
type Error struct {
	// HTTP status code of the error.
	StatusCode int

	// AWS code of the error.
	Code string

	// Message explaining the error.
	Message string
}

func (e *Error) Error() string {
	var prefix string
	if e.Code != "" {
		prefix = e.Code + ": "
	}
	if prefix == "" && e.StatusCode > 0 {
		prefix = strconv.Itoa(e.StatusCode) + ": "
	}
	return prefix + e.Message
}
