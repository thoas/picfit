package autoscaling_test

import (
	"testing"

	"github.com/mitchellh/goamz/autoscaling"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/testutil"
	. "github.com/motain/gocheck"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct {
	autoscaling *autoscaling.AutoScaling
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	auth := aws.Auth{"abc", "123", ""}
	s.autoscaling = autoscaling.NewWithClient(auth, aws.Region{AutoScalingEndpoint: testServer.URL}, testutil.DefaultClient)
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

func (s *S) Test_CreateAutoScalingGroup(c *C) {
	testServer.Response(200, nil, CreateAutoScalingGroupExample)

	options := autoscaling.CreateAutoScalingGroup{
		AvailZone:               []string{"us-east-1a"},
		DefaultCooldown:         30,
		DesiredCapacity:         2,
		HealthCheckGracePeriod:  30,
		HealthCheckType:         "elb",
		InstanceId:              "i-foo",
		LaunchConfigurationName: "foobar",
		MinSize:                 2,
		MaxSize:                 2,
		PlacementGroup:          "foobar",
		TerminationPolicies:     []string{"ClosestToNextInstanceHour", "OldestInstance"},
		Name:                    "foobar",
		Tags: []autoscaling.Tag{
			autoscaling.Tag{
				Key:   "foo",
				Value: "bar",
			},
		},
		VPCZoneIdentifier: []string{"foo", "bar"},
	}

	resp, err := s.autoscaling.CreateAutoScalingGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateAutoScalingGroup"})
	c.Assert(req.Form["InstanceId"], DeepEquals, []string{"i-foo"})
	c.Assert(req.Form["VPCZoneIdentifier"], DeepEquals, []string{"foo,bar"})
	c.Assert(req.Form["TerminationPolicies.member.1"], DeepEquals, []string{"ClosestToNextInstanceHour"})
	c.Assert(req.Form["TerminationPolicies.member.2"], DeepEquals, []string{"OldestInstance"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) Test_CreateLaunchConfiguration(c *C) {
	testServer.Response(200, nil, CreateLaunchConfigurationExample)

	options := autoscaling.CreateLaunchConfiguration{
		SecurityGroups: []string{"sg-1111"},
		ImageId:        "i-141421",
		InstanceId:     "i-141421",
		InstanceType:   "m1.small",
		KeyName:        "foobar",
		Name:           "i-141421",
		UserData:       "#!/bin/bash\necho Hello\n",
		BlockDevices:   []autoscaling.BlockDeviceMapping{
			{DeviceName: "/dev/sdb", VirtualName: "ephemeral0"},
			{DeviceName: "/dev/sdc", SnapshotId: "snap-a08912c9", DeleteOnTermination: true},
		},
	}

	resp, err := s.autoscaling.CreateLaunchConfiguration(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateLaunchConfiguration"})
	c.Assert(req.Form["InstanceType"], DeepEquals, []string{"m1.small"})
	c.Assert(req.Form["SecurityGroups.member.1"], DeepEquals, []string{"sg-1111"})
	c.Assert(req.Form["UserData"], DeepEquals, []string{"IyEvYmluL2Jhc2gKZWNobyBIZWxsbwo="})
	c.Assert(req.Form["BlockDeviceMappings.member.1.DeviceName"], DeepEquals, []string{"/dev/sdb"})
	c.Assert(req.Form["BlockDeviceMappings.member.1.VirtualName"], DeepEquals, []string{"ephemeral0"})
	c.Assert(req.Form["BlockDeviceMappings.member.2.Ebs.SnapshotId"], DeepEquals, []string{"snap-a08912c9"})
	c.Assert(req.Form["BlockDeviceMappings.member.2.Ebs.DeleteOnTermination"], DeepEquals, []string{"true"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7c6e177f-f082-11e1-ac58-3714bEXAMPLE")
}

func (s *S) Test_DescribeAutoScalingGroups(c *C) {
	testServer.Response(200, nil, DescribeAutoScalingGroupsExample)

	options := autoscaling.DescribeAutoScalingGroups{
		Names: []string{"foobar"},
	}

	resp, err := s.autoscaling.DescribeAutoScalingGroups(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeAutoScalingGroups"})
	c.Assert(req.Form["AutoScalingGroupNames.member.1"], DeepEquals, []string{"foobar"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "0f02a07d-b677-11e2-9eb0-dd50EXAMPLE")
	c.Assert(resp.AutoScalingGroups[0].Name, Equals, "my-test-asg-lbs")
	c.Assert(resp.AutoScalingGroups[0].LaunchConfigurationName, Equals, "my-test-lc")
	c.Assert(resp.AutoScalingGroups[0].TerminationPolicies[0], Equals, "Default")
}

func (s *S) Test_DescribeLaunchConfigurations(c *C) {
	testServer.Response(200, nil, DescribeLaunchConfigurationsExample)

	options := autoscaling.DescribeLaunchConfigurations{
		Names: []string{"foobar"},
	}

	resp, err := s.autoscaling.DescribeLaunchConfigurations(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DescribeLaunchConfigurations"})
	c.Assert(req.Form["LaunchConfigurationNames.member.1"], DeepEquals, []string{"foobar"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "d05a22f8-b690-11e2-bf8e-2113fEXAMPLE")
	c.Assert(resp.LaunchConfigurations[0].InstanceType, Equals, "m1.small")
}

func (s *S) TestDeleteAutoScalingGroup(c *C) {
	testServer.Response(200, nil, DeleteAutoScalingGroupExample)

	options := autoscaling.DeleteAutoScalingGroup{
		Name:        "foobar",
		ForceDelete: true,
	}

	resp, err := s.autoscaling.DeleteAutoScalingGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteAutoScalingGroup"})
	c.Assert(req.Form["AutoScalingGroupName"], DeepEquals, []string{"foobar"})
	c.Assert(req.Form["ForceDelete"], DeepEquals, []string{"true"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "70a76d42-9665-11e2-9fdf-211deEXAMPLE")
}

func (s *S) TestDeleteLaunchConfiguration(c *C) {
	testServer.Response(200, nil, DeleteLaunchConfigurationExample)

	options := autoscaling.DeleteLaunchConfiguration{
		Name: "foobar",
	}

	resp, err := s.autoscaling.DeleteLaunchConfiguration(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteLaunchConfiguration"})
	c.Assert(req.Form["LaunchConfigurationName"], DeepEquals, []string{"foobar"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7347261f-97df-11e2-8756-35eEXAMPLE")
}

func (s *S) Test_UpdateAutoScalingGroup(c *C) {
	testServer.Response(200, nil, UpdateAutoScalingGroupExample)

	options := autoscaling.UpdateAutoScalingGroup{
		AvailZone:       []string{"us-east-1a"},
		DefaultCooldown: 30,
		Name:            "bar",

		SetDefaultCooldown: true,
	}

	resp, err := s.autoscaling.UpdateAutoScalingGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"UpdateAutoScalingGroup"})
	c.Assert(req.Form["AutoScalingGroupName"], DeepEquals, []string{"bar"})
	c.Assert(req.Form["DefaultCooldown"], DeepEquals, []string{"30"})
	c.Assert(req.Form["MinSize"], IsNil)
	c.Assert(req.Form["MaxSize"], IsNil)
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "adafead0-ab8a-11e2-ba13-ab0ccEXAMPLE")
}
