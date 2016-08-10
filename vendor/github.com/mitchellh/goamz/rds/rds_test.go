package rds_test

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/rds"
	"github.com/mitchellh/goamz/testutil"
	. "github.com/motain/gocheck"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct {
	rds *rds.Rds
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	auth := aws.Auth{"abc", "123", ""}
	s.rds = rds.NewWithClient(auth, aws.Region{RdsEndpoint: testServer.URL}, testutil.DefaultClient)
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

func (s *S) Test_CreateDBInstance(c *C) {
	testServer.Response(200, nil, CreateDBInstanceExample)

	options := rds.CreateDBInstance{
		BackupRetentionPeriod:      30,
		MultiAZ:                    false,
		DBInstanceIdentifier:       "foobarbaz",
		PreferredBackupWindow:      "10:07-10:37",
		PreferredMaintenanceWindow: "sun:06:13-sun:06:43",
		AvailabilityZone:           "us-west-2b",
		Engine:                     "mysql",
		EngineVersion:              "",
		DBName:                     "5.6.13",
		AllocatedStorage:           10,
		StorageType:                "gp2",
		MasterUsername:             "foobar",
		MasterUserPassword:         "bazbarbaz",
		DBInstanceClass:            "db.m1.small",
		DBSecurityGroupNames:       []string{"foo", "bar"},
		DBParameterGroupName:       "default.mysql5.6",

		SetBackupRetentionPeriod: true,
	}

	resp, err := s.rds.CreateDBInstance(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateDBInstance"})
	c.Assert(req.Form["Engine"], DeepEquals, []string{"mysql"})
	c.Assert(req.Form["StorageType"], DeepEquals, []string{"gp2"})
	c.Assert(req.Form["DBSecurityGroups.member.1"], DeepEquals, []string{"foo"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "523e3218-afc7-11c3-90f5-f90431260ab4")
}

func (s *S) Test_CreateDBSecurityGroup(c *C) {
	testServer.Response(200, nil, CreateDBSecurityGroupExample)

	options := rds.CreateDBSecurityGroup{
		DBSecurityGroupName:        "foobarbaz",
		DBSecurityGroupDescription: "test description",
	}

	resp, err := s.rds.CreateDBSecurityGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateDBSecurityGroup"})
	c.Assert(req.Form["DBSecurityGroupName"], DeepEquals, []string{"foobarbaz"})
	c.Assert(req.Form["DBSecurityGroupDescription"], DeepEquals, []string{"test description"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "e68ef6fa-afc1-11c3-845a-476777009d19")
}

func (s *S) Test_CreateDBSubnetGroup(c *C) {
	testServer.Response(200, nil, CreateDBSubnetGroupExample)

	options := rds.CreateDBSubnetGroup{
		DBSubnetGroupName:        "foobarbaz",
		DBSubnetGroupDescription: "test description",
		SubnetIds:                []string{"subnet-e4d398a1", "subnet-c2bdb6ba"},
	}

	resp, err := s.rds.CreateDBSubnetGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateDBSubnetGroup"})
	c.Assert(req.Form["DBSubnetGroupName"], DeepEquals, []string{"foobarbaz"})
	c.Assert(req.Form["DBSubnetGroupDescription"], DeepEquals, []string{"test description"})
	c.Assert(req.Form["SubnetIds.member.1"], DeepEquals, []string{"subnet-e4d398a1"})
	c.Assert(req.Form["SubnetIds.member.2"], DeepEquals, []string{"subnet-c2bdb6ba"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "3a401b3f-bb9e-11d3-f4c6-37db295f7674")
}

func (s *S) Test_CreateDBParameterGroup(c *C) {
	testServer.Response(200, nil, CreateDBParameterGroupExample)

	options := rds.CreateDBParameterGroup{
		DBParameterGroupFamily: "mysql5.6",
		DBParameterGroupName:   "mydbparamgroup3",
		Description:            "My new DB Parameter Group",
	}

	resp, err := s.rds.CreateDBParameterGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateDBParameterGroup"})
	c.Assert(req.Form["DBParameterGroupFamily"], DeepEquals, []string{"mysql5.6"})
	c.Assert(req.Form["DBParameterGroupName"], DeepEquals, []string{"mydbparamgroup3"})
	c.Assert(req.Form["Description"], DeepEquals, []string{"My new DB Parameter Group"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7805c127-af22-11c3-96ac-6999cc5f7e72")
}

func (s *S) Test_DescribeDBInstances(c *C) {
	testServer.Response(200, nil, DescribeDBInstancesExample)

	options := rds.DescribeDBInstances{
		DBInstanceIdentifier: "foobarbaz",
	}

	resp, err := s.rds.DescribeDBInstances(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeDBInstances"})
	c.Assert(req.Form["DBInstanceIdentifier"], DeepEquals, []string{"foobarbaz"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "01b2685a-b978-11d3-f272-7cd6cce12cc5")
	c.Assert(resp.DBInstances[0].DBName, Equals, "mysampledb")
	c.Assert(resp.DBInstances[0].DBSecurityGroupNames, DeepEquals, []string{"my-db-secgroup"})
	c.Assert(resp.DBInstances[0].DBParameterGroupName, Equals, "default.mysql5.6")
	c.Assert(resp.DBInstances[0].StorageType, Equals, "gp2")
	c.Assert(resp.DBInstances[1].VpcSecurityGroupIds, DeepEquals, []string{"my-vpc-secgroup"})
}

func (s *S) Test_DescribeDBSecurityGroups(c *C) {
	testServer.Response(200, nil, DescribeDBSecurityGroupsExample)

	options := rds.DescribeDBSecurityGroups{
		DBSecurityGroupName: "foobarbaz",
	}

	resp, err := s.rds.DescribeDBSecurityGroups(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeDBSecurityGroups"})
	c.Assert(req.Form["DBSecurityGroupName"], DeepEquals, []string{"foobarbaz"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "b76e692c-b98c-11d3-a907-5a2c468b9cb0")
	c.Assert(resp.DBSecurityGroups[0].EC2SecurityGroupIds, DeepEquals, []string{"sg-7f476617"})
	c.Assert(resp.DBSecurityGroups[0].EC2SecurityGroupOwnerIds, DeepEquals, []string{"803#########"})
	c.Assert(resp.DBSecurityGroups[0].EC2SecurityGroupStatuses, DeepEquals, []string{"authorized"})
	c.Assert(resp.DBSecurityGroups[0].CidrIps, DeepEquals, []string{"192.0.0.0/24", "190.0.1.0/29", "190.0.2.0/29", "10.0.0.0/8"})
	c.Assert(resp.DBSecurityGroups[0].CidrStatuses, DeepEquals, []string{"authorized", "authorized", "authorized", "authorized"})
}

func (s *S) Test_DescribeDBSubnetGroups(c *C) {
	testServer.Response(200, nil, DescribeDBSubnetGroupsExample)

	options := rds.DescribeDBSubnetGroups{
		DBSubnetGroupName: "foobarbaz",
	}

	resp, err := s.rds.DescribeDBSubnetGroups(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeDBSubnetGroups"})
	c.Assert(req.Form["DBSubnetGroupName"], DeepEquals, []string{"foobarbaz"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "b783db3b-b98c-11d3-fbc7-5c0aad74da7c")
	c.Assert(resp.DBSubnetGroups[0].Status, DeepEquals, "Complete")
	c.Assert(resp.DBSubnetGroups[0].SubnetIds, DeepEquals, []string{"subnet-e8b3e5b1", "subnet-44b2f22e"})
	c.Assert(resp.DBSubnetGroups[0].VpcId, DeepEquals, "vpc-e7abbdce")
}

func (s *S) Test_DescribeDBParameterGroups(c *C) {
	testServer.Response(200, nil, DescribeDBParameterGroupsExample)

	options := rds.DescribeDBParameterGroups{
		DBParameterGroupName: "mydbparamgroup3",
	}

	resp, err := s.rds.DescribeDBParameterGroups(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeDBParameterGroups"})
	c.Assert(req.Form["DBParameterGroupName"], DeepEquals, []string{"mydbparamgroup3"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "b75d527a-b98c-11d3-f272-7cd6cce12cc5")
	c.Assert(resp.DBParameterGroups[0].DBParameterGroupFamily, Equals, "mysql5.6")
	c.Assert(resp.DBParameterGroups[0].Description, Equals, "My new DB Parameter Group")
	c.Assert(resp.DBParameterGroups[0].DBParameterGroupName, Equals, "mydbparamgroup3")
}

func (s *S) Test_DescribeDBParameters(c *C) {
	testServer.Response(200, nil, DescribeDBParametersExample)

	options := rds.DescribeDBParameters{
		DBParameterGroupName: "mydbparamgroup3",
		Source:               "user",
	}

	resp, err := s.rds.DescribeDBParameters(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeDBParameters"})
	c.Assert(req.Form["DBParameterGroupName"], DeepEquals, []string{"mydbparamgroup3"})
	c.Assert(req.Form["Source"], DeepEquals, []string{"user"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "8c40488f-b9ff-11d3-a15e-7ac49293f4fa")
	c.Assert(resp.Parameters[0].ParameterName, Equals, "character_set_server")
	c.Assert(resp.Parameters[0].ParameterValue, Equals, "utf8")
	c.Assert(resp.Parameters[1].ParameterName, Equals, "character_set_client")
	c.Assert(resp.Parameters[1].ParameterValue, Equals, "utf8")
	c.Assert(resp.Parameters[2].ParameterName, Equals, "character_set_results")
	c.Assert(resp.Parameters[2].ParameterValue, Equals, "utf8")
	c.Assert(resp.Parameters[3].ParameterName, Equals, "collation_server")
	c.Assert(resp.Parameters[3].ParameterValue, Equals, "utf8_unicode_ci")
	c.Assert(resp.Parameters[4].ParameterName, Equals, "collation_connection")
	c.Assert(resp.Parameters[4].ParameterValue, Equals, "utf8_unicode_ci")
}

func (s *S) Test_DeleteDBInstance(c *C) {
	testServer.Response(200, nil, DeleteDBInstanceExample)

	options := rds.DeleteDBInstance{
		DBInstanceIdentifier: "foobarbaz",
		SkipFinalSnapshot:    true,
	}

	resp, err := s.rds.DeleteDBInstance(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteDBInstance"})
	c.Assert(req.Form["DBInstanceIdentifier"], DeepEquals, []string{"foobarbaz"})
	c.Assert(req.Form["SkipFinalSnapshot"], DeepEquals, []string{"true"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7369556f-b70d-11c3-faca-6ba18376ea1b")
}

func (s *S) Test_DeleteDBInstance_SnapshotIdentifier(c *C) {
	testServer.Response(200, nil, DeleteDBInstanceExample)

	options := rds.DeleteDBInstance{
		DBInstanceIdentifier:      "foobarbaz",
		SkipFinalSnapshot:         false,
		FinalDBSnapshotIdentifier: "bar",
	}

	resp, err := s.rds.DeleteDBInstance(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteDBInstance"})
	c.Assert(req.Form["DBInstanceIdentifier"], DeepEquals, []string{"foobarbaz"})
	c.Assert(req.Form["FinalDBSnapshotIdentifier"], DeepEquals, []string{"bar"})
	c.Assert(req.Form["SkipFinalSnapshot"], IsNil)
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7369556f-b70d-11c3-faca-6ba18376ea1b")
}

func (s *S) Test_DeleteDBSecurityGroup(c *C) {
	testServer.Response(200, nil, DeleteDBSecurityGroupExample)

	options := rds.DeleteDBSecurityGroup{
		DBSecurityGroupName: "foobarbaz",
	}

	resp, err := s.rds.DeleteDBSecurityGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteDBSecurityGroup"})
	c.Assert(req.Form["DBSecurityGroupName"], DeepEquals, []string{"foobarbaz"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7aec7454-ba25-11d3-855b-576787000e19")
}

func (s *S) Test_DeleteDBSubnetGroup(c *C) {
	testServer.Response(200, nil, DeleteDBSubnetGroupExample)

	options := rds.DeleteDBSubnetGroup{
		DBSubnetGroupName: "foobarbaz",
	}

	resp, err := s.rds.DeleteDBSubnetGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteDBSubnetGroup"})
	c.Assert(req.Form["DBSubnetGroupName"], DeepEquals, []string{"foobarbaz"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "6295e5ab-bbf3-11d3-f4c6-37db295f7674")
}

func (s *S) Test_DeleteDBParameterGroup(c *C) {
	testServer.Response(200, nil, DeleteDBParameterGroupExample)

	options := rds.DeleteDBParameterGroup{
		DBParameterGroupName: "mydbparamgroup3",
	}

	resp, err := s.rds.DeleteDBParameterGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteDBParameterGroup"})
	c.Assert(req.Form["DBParameterGroupName"], DeepEquals, []string{"mydbparamgroup3"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "cad6c267-ba25-11d3-fe11-33d33a9bb7e3")
}

func (s *S) Test_AuthorizeDBSecurityGroupIngress(c *C) {
	testServer.Response(200, nil, AuthorizeDBSecurityGroupIngressExample)

	options := rds.AuthorizeDBSecurityGroupIngress{
		DBSecurityGroupName:     "foobarbaz",
		EC2SecurityGroupOwnerId: "bar",
	}

	resp, err := s.rds.AuthorizeDBSecurityGroupIngress(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"AuthorizeDBSecurityGroupIngress"})
	c.Assert(req.Form["DBSecurityGroupName"], DeepEquals, []string{"foobarbaz"})
	c.Assert(req.Form["EC2SecurityGroupOwnerId"], DeepEquals, []string{"bar"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "6176b5f8-bfed-11d3-f92b-31fa5e8dbc99")
}

func (s *S) Test_DescribeDBSnapshots(c *C) {
	testServer.Response(200, nil, DescribeDBSnapshotsExample)

	options := rds.DescribeDBSnapshots{
		DBInstanceIdentifier: "foobar",
		DBSnapshotIdentifier: "baz",
		SnapshotType:         "manual",
	}

	resp, err := s.rds.DescribeDBSnapshots(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeDBSnapshots"})
	c.Assert(req.Form["DBInstanceIdentifier"], DeepEquals, []string{"foobar"})
	c.Assert(req.Form["DBSnapshotIdentifier"], DeepEquals, []string{"baz"})
	c.Assert(req.Form["SnapshotType"], DeepEquals, []string{"manual"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "b7769930-b98c-11d3-f272-7cd6cce12cc5")
	c.Assert(resp.DBSnapshots[0].OptionGroupName, Equals, "default:mysql-5-6")
	c.Assert(resp.DBSnapshots[0].Engine, Equals, "mysql")
	c.Assert(resp.DBSnapshots[0].SnapshotType, Equals, "manual")
}

func (s *S) Test_RestoreDBInstanceFromDBSnapshot(c *C) {
	testServer.Response(200, nil, RestoreDBInstanceFromDBSnapshotExample)

	options := rds.RestoreDBInstanceFromDBSnapshot{
		DBInstanceIdentifier: "foo",
		DBSnapshotIdentifier: "bar",
	}

	resp, err := s.rds.RestoreDBInstanceFromDBSnapshot(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"RestoreDBInstanceFromDBSnapshot"})
	c.Assert(req.Form["DBInstanceIdentifier"], DeepEquals, []string{"foo"})
	c.Assert(req.Form["DBSnapshotIdentifier"], DeepEquals, []string{"bar"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "863fd73e-be2b-11d3-855b-576787000e19")
}

func (s *S) Test_ModifyDBParameterGroup(c *C) {
	testServer.Response(200, nil, ModifyDBParameterGroupExample)

	options := rds.ModifyDBParameterGroup{
		DBParameterGroupName: "mydbparamgroup3",
		Parameters: []rds.Parameter{
			rds.Parameter{
				ApplyMethod:    "immediate",
				ParameterName:  "character_set_server",
				ParameterValue: "utf8",
			},
			rds.Parameter{
				ApplyMethod:    "immediate",
				ParameterName:  "character_set_client",
				ParameterValue: "utf8",
			},
			rds.Parameter{
				ApplyMethod:    "immediate",
				ParameterName:  "character_set_results",
				ParameterValue: "utf8",
			},
			rds.Parameter{
				ApplyMethod:    "immediate",
				ParameterName:  "collation_server",
				ParameterValue: "utf8_unicode_ci",
			},
			rds.Parameter{
				ApplyMethod:    "immediate",
				ParameterName:  "collation_connection",
				ParameterValue: "utf8_unicode_ci",
			},
		},
	}

	resp, err := s.rds.ModifyDBParameterGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"ModifyDBParameterGroup"})
	c.Assert(req.Form["DBParameterGroupName"], DeepEquals, []string{"mydbparamgroup3"})
	c.Assert(req.Form["Parameters.member.1.ApplyMethod"], DeepEquals, []string{"immediate"})
	c.Assert(req.Form["Parameters.member.1.ParameterName"], DeepEquals, []string{"character_set_server"})
	c.Assert(req.Form["Parameters.member.1.ParameterValue"], DeepEquals, []string{"utf8"})
	c.Assert(req.Form["Parameters.member.2.ApplyMethod"], DeepEquals, []string{"immediate"})
	c.Assert(req.Form["Parameters.member.2.ParameterName"], DeepEquals, []string{"character_set_client"})
	c.Assert(req.Form["Parameters.member.2.ParameterValue"], DeepEquals, []string{"utf8"})
	c.Assert(req.Form["Parameters.member.3.ApplyMethod"], DeepEquals, []string{"immediate"})
	c.Assert(req.Form["Parameters.member.3.ParameterName"], DeepEquals, []string{"character_set_results"})
	c.Assert(req.Form["Parameters.member.3.ParameterValue"], DeepEquals, []string{"utf8"})
	c.Assert(req.Form["Parameters.member.4.ApplyMethod"], DeepEquals, []string{"immediate"})
	c.Assert(req.Form["Parameters.member.4.ParameterName"], DeepEquals, []string{"collation_server"})
	c.Assert(req.Form["Parameters.member.4.ParameterValue"], DeepEquals, []string{"utf8_unicode_ci"})
	c.Assert(req.Form["Parameters.member.5.ApplyMethod"], DeepEquals, []string{"immediate"})
	c.Assert(req.Form["Parameters.member.5.ParameterName"], DeepEquals, []string{"collation_connection"})
	c.Assert(req.Form["Parameters.member.5.ParameterValue"], DeepEquals, []string{"utf8_unicode_ci"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "12d7435e-bba0-11d3-fe11-33d33a9bb7e3")
}
