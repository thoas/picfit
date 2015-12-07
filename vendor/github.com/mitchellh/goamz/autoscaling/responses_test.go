package autoscaling_test

var ErrorDump = `
<?xml version="1.0" encoding="UTF-8"?>
<Response><Errors><Error><Code>UnsupportedOperation</Code>
<Message></Message>
</Error></Errors><RequestID>0503f4e9-bbd6-483c-b54f-c4ae9f3b30f4</RequestID></Response>
`

// http://goo.gl/gQRD2H
var CreateAutoScalingGroupExample = `
<CreateAutoScalingGroupResponse xmlns="http://autoscaling.amazonaws.com/doc/2011-01-01/">
<ResponseMetadata>
<RequestId>8d798a29-f083-11e1-bdfb-cb223EXAMPLE</RequestId>
</ResponseMetadata>
</CreateAutoScalingGroupResponse>
`

var CreateLaunchConfigurationExample = `
<CreateLaunchConfigurationResponse xmlns="http://autoscaling.amazonaws.com/doc/2011-01-01/">
<ResponseMetadata>
   <RequestId>7c6e177f-f082-11e1-ac58-3714bEXAMPLE</RequestId>
</ResponseMetadata>
</CreateLaunchConfigurationResponse>
`

var DescribeLaunchConfigurationsExample = `
<DescribeLaunchConfigurationsResponse xmlns="http://autoscaling.amazonaws.com/doc/2011-01-01/">
  <DescribeLaunchConfigurationsResult>
    <LaunchConfigurations>
      <member>
        <AssociatePublicIpAddress>true</AssociatePublicIpAddress>
        <SecurityGroups/>
        <PlacementTenancy>dedicated</PlacementTenancy>
        <CreatedTime>2013-01-21T23:04:42.200Z</CreatedTime>
        <KernelId/>
        <LaunchConfigurationName>my-test-lc</LaunchConfigurationName>
        <UserData/>
        <InstanceType>m1.small</InstanceType>
        <LaunchConfigurationARN>arn:aws:autoscaling:us-east-1:803981987763:launchConfiguration:
        9dbbbf87-6141-428a-a409-0752edbe6cad:launchConfigurationName/my-test-lc</LaunchConfigurationARN>
        <BlockDeviceMappings/>
        <ImageId>ami-514ac838</ImageId>
        <KeyName/>
        <RamdiskId/>
        <InstanceMonitoring>
          <Enabled>true</Enabled>
        </InstanceMonitoring>
        <EbsOptimized>false</EbsOptimized>
      </member>
    </LaunchConfigurations>
  </DescribeLaunchConfigurationsResult>
  <ResponseMetadata>
    <RequestId>d05a22f8-b690-11e2-bf8e-2113fEXAMPLE</RequestId>
  </ResponseMetadata>
</DescribeLaunchConfigurationsResponse>
`

var DescribeAutoScalingGroupsExample = `
<DescribeAutoScalingGroupsResponse xmlns="http://autoscaling.amazonaws.com/doc/2011-01-01/">
<DescribeAutoScalingGroupsResult>
    <AutoScalingGroups>
      <member>
        <Tags/>
        <SuspendedProcesses/>
        <AutoScalingGroupName>my-test-asg-lbs</AutoScalingGroupName>
        <HealthCheckType>ELB</HealthCheckType>
        <CreatedTime>2013-05-06T17:47:15.107Z</CreatedTime>
        <EnabledMetrics/>
        <LaunchConfigurationName>my-test-lc</LaunchConfigurationName>
        <Instances/>
        <DesiredCapacity>2</DesiredCapacity>
        <AvailabilityZones>
          <member>us-east-1b</member>
          <member>us-east-1a</member>
        </AvailabilityZones>
        <LoadBalancerNames>
          <member>my-test-asg-loadbalancer</member>
        </LoadBalancerNames>
        <MinSize>2</MinSize>
        <VPCZoneIdentifier/>
        <HealthCheckGracePeriod>120</HealthCheckGracePeriod>
        <DefaultCooldown>300</DefaultCooldown>
        <AutoScalingGroupARN>arn:aws:autoscaling:us-east-1:803981987763:autoScalingGroup:ca861182-c8f9-4ca7-b1eb-cd35505f5ebb
        :autoScalingGroupName/my-test-asg-lbs</AutoScalingGroupARN>
        <TerminationPolicies>
          <member>Default</member>
        </TerminationPolicies>
        <MaxSize>10</MaxSize>
      </member>
    </AutoScalingGroups>
  </DescribeAutoScalingGroupsResult>
  <ResponseMetadata>
    <RequestId>0f02a07d-b677-11e2-9eb0-dd50EXAMPLE</RequestId>
  </ResponseMetadata>
</DescribeAutoScalingGroupsResponse>
`

var DeleteLaunchConfigurationExample = `<DeleteLaunchConfigurationResponse xmlns="http://autoscaling.amazonaws.com/doc/2011-01-01/">
  <ResponseMetadata>
    <RequestId>7347261f-97df-11e2-8756-35eEXAMPLE</RequestId>
  </ResponseMetadata>
</DeleteLaunchConfigurationResponse>`

var DeleteAutoScalingGroupExample = `<DeleteAutoScalingGroupResponse xmlns="http://autoscaling.amazonaws.com/doc/2011-01-01/">
  <ResponseMetadata>
    <RequestId>70a76d42-9665-11e2-9fdf-211deEXAMPLE</RequestId>
  </ResponseMetadata>
</DeleteAutoScalingGroupResponse>`

var UpdateAutoScalingGroupExample = `
<UpdateAutoScalingGroupResponse xmlns="http://autoscaling.amazonaws.com/doc/2011-01-01/">
  <ResponseMetadata>
    <RequestId>adafead0-ab8a-11e2-ba13-ab0ccEXAMPLE</RequestId>
  </ResponseMetadata>
</UpdateAutoScalingGroupResponse>
`
