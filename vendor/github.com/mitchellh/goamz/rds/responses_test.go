package rds_test

var ErrorDump = `
<?xml version="1.0" encoding="UTF-8"?>
<Response><Errors><Error><Code>UnsupportedOperation</Code>
<Message></Message>
</Error></Errors><RequestID>0503f4e9-bbd6-483c-b54f-c4ae9f3b30f4</RequestID></Response>
`

// http://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_DescribeDBInstances.html
var DescribeDBInstancesExample = `
<DescribeDBInstancesResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <DescribeDBInstancesResult>
    <DBInstances>
      <DBInstance>
        <BackupRetentionPeriod>7</BackupRetentionPeriod>
        <MultiAZ>false</MultiAZ>
        <DBInstanceStatus>available</DBInstanceStatus>
        <VpcSecurityGroups/>
        <DBInstanceIdentifier>mysqlexampledb</DBInstanceIdentifier>
        <PreferredBackupWindow>10:07-10:37</PreferredBackupWindow>
        <PreferredMaintenanceWindow>sun:06:13-sun:06:43</PreferredMaintenanceWindow>
        <AvailabilityZone>us-west-2b</AvailabilityZone>
        <LatestRestorableTime>2014-04-21T17:15:00Z</LatestRestorableTime>
        <ReadReplicaDBInstanceIdentifiers/>
        <Engine>mysql</Engine>
        <PendingModifiedValues/>
        <LicenseModel>general-public-license</LicenseModel>
        <DBParameterGroups>
          <DBParameterGroup>
            <ParameterApplyStatus>in-sync</ParameterApplyStatus>
            <DBParameterGroupName>default.mysql5.6</DBParameterGroupName>
          </DBParameterGroup>
        </DBParameterGroups>
        <Endpoint>
          <Port>3306</Port>
          <Address>mysqlexampledb.c6c1rntzufv0.us-west-2.rds.amazonaws.com</Address>
        </Endpoint>
        <EngineVersion>5.6.13</EngineVersion>
        <OptionGroupMemberships>
          <OptionGroupMembership>
            <OptionGroupName>default:mysql-5-6</OptionGroupName>
            <Status>in-sync</Status>
          </OptionGroupMembership>
        </OptionGroupMemberships>
        <DBSecurityGroups>
          <DBSecurityGroup>
            <Status>active</Status>
            <DBSecurityGroupName>my-db-secgroup</DBSecurityGroupName>
          </DBSecurityGroup>
        </DBSecurityGroups>
        <PubliclyAccessible>true</PubliclyAccessible>
        <DBName>mysampledb</DBName>
        <AutoMinorVersionUpgrade>true</AutoMinorVersionUpgrade>
        <InstanceCreateTime>2014-01-29T22:58:24.231Z</InstanceCreateTime>
        <AllocatedStorage>5</AllocatedStorage>
        <StorageType>gp2</StorageType>
        <MasterUsername>myawsuser</MasterUsername>
        <DBInstanceClass>db.t1.micro</DBInstanceClass>
      </DBInstance>
      <DBInstance>
        <BackupRetentionPeriod>7</BackupRetentionPeriod>
        <MultiAZ>false</MultiAZ>
        <DBInstanceStatus>available</DBInstanceStatus>
        <VpcSecurityGroups>
          <VpcSecurityGroupMembership>
            <Status>active</Status>
            <VpcSecurityGroupId>my-vpc-secgroup</VpcSecurityGroupId>
          </VpcSecurityGroupMembership>
        </VpcSecurityGroups>
        <DBInstanceIdentifier>mysqlexampledb-restore</DBInstanceIdentifier>
        <PreferredBackupWindow>10:07-10:37</PreferredBackupWindow>
        <PreferredMaintenanceWindow>sun:06:13-sun:06:43</PreferredMaintenanceWindow>
        <AvailabilityZone>us-west-2b</AvailabilityZone>
        <LatestRestorableTime>2014-04-21T17:15:00Z</LatestRestorableTime>
        <ReadReplicaDBInstanceIdentifiers/>
        <Engine>mysql</Engine>
        <PendingModifiedValues/>
        <LicenseModel>general-public-license</LicenseModel>
        <DBParameterGroups>
          <DBParameterGroup>
            <ParameterApplyStatus>in-sync</ParameterApplyStatus>
            <DBParameterGroupName>default.mysql5.6</DBParameterGroupName>
          </DBParameterGroup>
        </DBParameterGroups>
        <Endpoint>
          <Port>3306</Port>
          <Address>mysqlexampledb-restore.c6c2mntzugv0.us-west-2.rds.amazonaws.com</Address>
        </Endpoint>
        <EngineVersion>5.6.13</EngineVersion>
        <OptionGroupMemberships>
          <OptionGroupMembership>
            <OptionGroupName>default:mysql-5-6</OptionGroupName>
            <Status>in-sync</Status>
          </OptionGroupMembership>
        </OptionGroupMemberships>
        <DBSecurityGroups/>
        <PubliclyAccessible>true</PubliclyAccessible>
        <DBName>mysampledb</DBName>
        <AutoMinorVersionUpgrade>true</AutoMinorVersionUpgrade>
        <InstanceCreateTime>2014-03-28T20:14:17.296Z</InstanceCreateTime>
        <AllocatedStorage>5</AllocatedStorage>
        <MasterUsername>myawsuser</MasterUsername>
        <DBInstanceClass>db.t1.micro</DBInstanceClass>
      </DBInstance>
    </DBInstances>
  </DescribeDBInstancesResult>
  <ResponseMetadata>
    <RequestId>01b2685a-b978-11d3-f272-7cd6cce12cc5</RequestId>
  </ResponseMetadata>
</DescribeDBInstancesResponse>
`

var CreateDBInstanceExample = `
<CreateDBInstanceResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <CreateDBInstanceResult>
    <DBInstance>
      <BackupRetentionPeriod>7</BackupRetentionPeriod>
      <DBInstanceStatus>creating</DBInstanceStatus>
      <MultiAZ>false</MultiAZ>
      <VpcSecurityGroups/>
      <DBInstanceIdentifier>myawsuser-dbi01</DBInstanceIdentifier>

      <PreferredBackupWindow>03:50-04:20</PreferredBackupWindow>
      <PreferredMaintenanceWindow>wed:06:38-wed:07:08</PreferredMaintenanceWindow>
      <ReadReplicaDBInstanceIdentifiers/>
      <Engine>mysql</Engine>
      <PendingModifiedValues>
        <MasterUserPassword>****</MasterUserPassword>
      </PendingModifiedValues>
      <LicenseModel>general-public-license</LicenseModel>
      <EngineVersion>5.6.13</EngineVersion>
      <DBParameterGroups>
        <DBParameterGroup>
          <ParameterApplyStatus>in-sync</ParameterApplyStatus>
          <DBParameterGroupName>default.mysql5.6</DBParameterGroupName>
        </DBParameterGroup>
      </DBParameterGroups>
      <OptionGroupMemberships>
        <OptionGroupMembership>
          <OptionGroupName>default:mysql-5-6</OptionGroupName>
          <Status>in-sync</Status>
        </OptionGroupMembership>
      </OptionGroupMemberships>
      <DBSecurityGroups>
        <DBSecurityGroup>
          <Status>active</Status>
          <DBSecurityGroupName>default</DBSecurityGroupName>
        </DBSecurityGroup>
      </DBSecurityGroups>
      <PubliclyAccessible>true</PubliclyAccessible>
      <AutoMinorVersionUpgrade>true</AutoMinorVersionUpgrade>
      <AllocatedStorage>15</AllocatedStorage>
      <DBInstanceClass>db.m1.large</DBInstanceClass>
      <MasterUsername>myawsuser</MasterUsername>
    </DBInstance>
  </CreateDBInstanceResult>
  <ResponseMetadata>
    <RequestId>523e3218-afc7-11c3-90f5-f90431260ab4</RequestId>
  </ResponseMetadata>
</CreateDBInstanceResponse>
`

var DeleteDBInstanceExample = `
<DeleteDBInstanceResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <DeleteDBInstanceResult>
    <DBInstance>
      <BackupRetentionPeriod>7</BackupRetentionPeriod>
      <DBInstanceStatus>deleting</DBInstanceStatus>
      <MultiAZ>false</MultiAZ>
      <VpcSecurityGroups/>
      <DBInstanceIdentifier>mydatabase</DBInstanceIdentifier>
      <PreferredBackupWindow>08:14-08:44</PreferredBackupWindow>
      <PreferredMaintenanceWindow>fri:04:50-fri:05:20</PreferredMaintenanceWindow>
      <AvailabilityZone>us-east-1a</AvailabilityZone>
      <ReadReplicaDBInstanceIdentifiers/>
      <LatestRestorableTime>2013-11-09T00:15:00Z</LatestRestorableTime>
      <Engine>mysql</Engine>
      <PendingModifiedValues/>
      <LicenseModel>general-public-license</LicenseModel>
      <EngineVersion>5.6.13</EngineVersion>
      <Endpoint>
        <Port>3306</Port>
        <Address>mydatabase.cf037hpkuvjt.us-east-1.rds.amazonaws.com</Address>
      </Endpoint>
      <DBParameterGroups>
        <DBParameterGroup>
          <ParameterApplyStatus>in-sync</ParameterApplyStatus>
          <DBParameterGroupName>default.mysql5.6</DBParameterGroupName>
        </DBParameterGroup>
      </DBParameterGroups>
      <OptionGroupMemberships>
        <OptionGroupMembership>
          <OptionGroupName>default:mysql-5-6</OptionGroupName>
          <Status>in-sync</Status>
        </OptionGroupMembership>
      </OptionGroupMemberships>
      <PubliclyAccessible>true</PubliclyAccessible>
      <DBSecurityGroups>
        <DBSecurityGroup>
          <Status>active</Status>
          <DBSecurityGroupName>default</DBSecurityGroupName>
        </DBSecurityGroup>
      </DBSecurityGroups>
      <DBName>mysqldb</DBName>
      <AutoMinorVersionUpgrade>true</AutoMinorVersionUpgrade>
      <InstanceCreateTime>2011-04-28T23:33:54.909Z</InstanceCreateTime>
      <AllocatedStorage>100</AllocatedStorage>
      <MasterUsername>myawsuser</MasterUsername>
      <DBInstanceClass>db.m1.medium</DBInstanceClass>
    </DBInstance>
  </DeleteDBInstanceResult>
  <ResponseMetadata>
    <RequestId>7369556f-b70d-11c3-faca-6ba18376ea1b</RequestId>
  </ResponseMetadata>
</DeleteDBInstanceResponse>
`

var DescribeDBSecurityGroupsExample = `
<DescribeDBSecurityGroupsResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <DescribeDBSecurityGroupsResult>
    <DBSecurityGroups>
      <DBSecurityGroup>
        <EC2SecurityGroups>
          <EC2SecurityGroup>
            <Status>authorized</Status>
            <EC2SecurityGroupName>elasticbeanstalk-windows</EC2SecurityGroupName>
            <EC2SecurityGroupOwnerId>803#########</EC2SecurityGroupOwnerId>
            <EC2SecurityGroupId>sg-7f476617</EC2SecurityGroupId>
          </EC2SecurityGroup>
        </EC2SecurityGroups>
        <DBSecurityGroupDescription>My security group</DBSecurityGroupDescription>
        <IPRanges>
          <IPRange>
            <CIDRIP>192.0.0.0/24</CIDRIP>
            <Status>authorized</Status>
          </IPRange>
          <IPRange>
            <CIDRIP>190.0.1.0/29</CIDRIP>
            <Status>authorized</Status>
          </IPRange>
          <IPRange>
            <CIDRIP>190.0.2.0/29</CIDRIP>
            <Status>authorized</Status>
          </IPRange>
          <IPRange>
            <CIDRIP>10.0.0.0/8</CIDRIP>
            <Status>authorized</Status>
          </IPRange>
        </IPRanges>
        <OwnerId>803#########</OwnerId>
        <DBSecurityGroupName>my-secgrp</DBSecurityGroupName>
      </DBSecurityGroup>
      <DBSecurityGroup>
        <EC2SecurityGroups/>
        <DBSecurityGroupDescription>default</DBSecurityGroupDescription>
        <IPRanges/>
        <OwnerId>803#########</OwnerId>
        <DBSecurityGroupName>default</DBSecurityGroupName>
      </DBSecurityGroup>
   </DBSecurityGroups>
  </DescribeDBSecurityGroupsResult>
  <ResponseMetadata>
    <RequestId>b76e692c-b98c-11d3-a907-5a2c468b9cb0</RequestId>
  </ResponseMetadata>
</DescribeDBSecurityGroupsResponse>
`

var DeleteDBSecurityGroupExample = `
<DeleteDBSecurityGroupResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <ResponseMetadata>
    <RequestId>7aec7454-ba25-11d3-855b-576787000e19</RequestId>
  </ResponseMetadata>
</DeleteDBSecurityGroupResponse>
`

var CreateDBSecurityGroupExample = `
<CreateDBSecurityGroupResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <CreateDBSecurityGroupResult>
    <DBSecurityGroup>
      <EC2SecurityGroups/>
      <DBSecurityGroupDescription>My new DB Security Group</DBSecurityGroupDescription>
      <IPRanges/>
      <OwnerId>803#########</OwnerId>
      <DBSecurityGroupName>mydbsecuritygroup00</DBSecurityGroupName>
    </DBSecurityGroup>
  </CreateDBSecurityGroupResult>
  <ResponseMetadata>
    <RequestId>e68ef6fa-afc1-11c3-845a-476777009d19</RequestId>
  </ResponseMetadata>
</CreateDBSecurityGroupResponse>
`

var AuthorizeDBSecurityGroupIngressExample = `
<AuthorizeDBSecurityGroupIngressResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <AuthorizeDBSecurityGroupIngressResult>
    <DBSecurityGroup>
      <EC2SecurityGroups>
        <EC2SecurityGroup>
          <Status>authorized</Status>
          <EC2SecurityGroupName>elasticbeanstalk-windows</EC2SecurityGroupName>
          <EC2SecurityGroupOwnerId>803#########</EC2SecurityGroupOwnerId>
          <EC2SecurityGroupId>sg-7f476617</EC2SecurityGroupId>
        </EC2SecurityGroup>
      </EC2SecurityGroups>
      <DBSecurityGroupDescription>default</DBSecurityGroupDescription>
      <IPRanges>
        <IPRange>
          <CIDRIP>192.0.0.0/24</CIDRIP>
          <Status>authorized</Status>
        </IPRange>
        <IPRange>
          <CIDRIP>190.0.1.0/29</CIDRIP>
          <Status>authorized</Status>
        </IPRange>
        <IPRange>
          <CIDRIP>190.0.2.0/29</CIDRIP>
          <Status>authorized</Status>
        </IPRange>
        <IPRange>
          <CIDRIP>10.0.0.0/8</CIDRIP>
          <Status>authorized</Status>
        </IPRange>
      </IPRanges>
      <OwnerId>803#########</OwnerId>
      <DBSecurityGroupName>default</DBSecurityGroupName>
    </DBSecurityGroup>
  </AuthorizeDBSecurityGroupIngressResult>
  <ResponseMetadata>
    <RequestId>6176b5f8-bfed-11d3-f92b-31fa5e8dbc99</RequestId>
  </ResponseMetadata>
</AuthorizeDBSecurityGroupIngressResponse>
`

var DescribeDBSubnetGroupsExample = `
<DescribeDBSubnetGroupsResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <DescribeDBSubnetGroupsResult>
    <DBSubnetGroups>
      <DBSubnetGroup>
        <VpcId>vpc-e7abbdce</VpcId>
        <SubnetGroupStatus>Complete</SubnetGroupStatus>
        <DBSubnetGroupDescription>DB subnet group 1</DBSubnetGroupDescription>
        <DBSubnetGroupName>mydbsubnetgroup1</DBSubnetGroupName>
        <Subnets>
          <Subnet>
            <SubnetStatus>Active</SubnetStatus>
            <SubnetIdentifier>subnet-e8b3e5b1</SubnetIdentifier>
            <SubnetAvailabilityZone>
              <Name>us-west-2a</Name>
              <ProvisionedIopsCapable>false</ProvisionedIopsCapable>
            </SubnetAvailabilityZone>
          </Subnet>
          <Subnet>
            <SubnetStatus>Active</SubnetStatus>
            <SubnetIdentifier>subnet-44b2f22e</SubnetIdentifier>
            <SubnetAvailabilityZone>
              <Name>us-west-2b</Name>
              <ProvisionedIopsCapable>false</ProvisionedIopsCapable>
            </SubnetAvailabilityZone>
          </Subnet>
        </Subnets>
      </DBSubnetGroup>
      <DBSubnetGroup>
        <VpcId>vpc-c1e17bb8</VpcId>
        <SubnetGroupStatus>Complete</SubnetGroupStatus>
        <DBSubnetGroupDescription>My DB Subnet Group 2</DBSubnetGroupDescription>
        <DBSubnetGroupName>sub-grp-2</DBSubnetGroupName>
        <Subnets>
          <Subnet>
            <SubnetStatus>Active</SubnetStatus>
            <SubnetIdentifier>subnet-d281ef8a</SubnetIdentifier>
            <SubnetAvailabilityZone>
              <Name>us-west-2a</Name>
              <ProvisionedIopsCapable>false</ProvisionedIopsCapable>
            </SubnetAvailabilityZone>
          </Subnet>
          <Subnet>
            <SubnetStatus>Active</SubnetStatus>
            <SubnetIdentifier>subnet-b381ef9f</SubnetIdentifier>
            <SubnetAvailabilityZone>
              <Name>us-west-2c</Name>
              <ProvisionedIopsCapable>false</ProvisionedIopsCapable>
            </SubnetAvailabilityZone>
          </Subnet>
          <Subnet>
            <SubnetStatus>Active</SubnetStatus>
            <SubnetIdentifier>subnet-e1e17ebd</SubnetIdentifier>
            <SubnetAvailabilityZone>
              <Name>us-west-2b</Name>
              <ProvisionedIopsCapable>false</ProvisionedIopsCapable>
            </SubnetAvailabilityZone>
          </Subnet>
        </Subnets>
      </DBSubnetGroup>
    </DBSubnetGroups>
  </DescribeDBSubnetGroupsResult>
  <ResponseMetadata>
    <RequestId>b783db3b-b98c-11d3-fbc7-5c0aad74da7c</RequestId>
  </ResponseMetadata>
</DescribeDBSubnetGroupsResponse>
`

var DeleteDBSubnetGroupExample = `
<DeleteDBSubnetGroupResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <ResponseMetadata>
    <RequestId>6295e5ab-bbf3-11d3-f4c6-37db295f7674</RequestId>
  </ResponseMetadata>
</DeleteDBSubnetGroupResponse>
`

var CreateDBSubnetGroupExample = `
<CreateDBSubnetGroupResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <CreateDBSubnetGroupResult>
    <DBSubnetGroup>
      <VpcId>vpc-33dc97ea</VpcId>
      <SubnetGroupStatus>Complete</SubnetGroupStatus>
      <DBSubnetGroupDescription>My new DB Subnet Group</DBSubnetGroupDescription>
      <DBSubnetGroupName>myawsuser-dbsubnetgroup</DBSubnetGroupName>
      <Subnets>
        <Subnet>
          <SubnetStatus>Active</SubnetStatus>
          <SubnetIdentifier>subnet-e4d398a1</SubnetIdentifier>
          <SubnetAvailabilityZone>
            <Name>us-east-1b</Name>
            <ProvisionedIopsCapable>false</ProvisionedIopsCapable>
          </SubnetAvailabilityZone>
        </Subnet>
        <Subnet>
          <SubnetStatus>Active</SubnetStatus>
          <SubnetIdentifier>subnet-c2bdb6ba</SubnetIdentifier>
          <SubnetAvailabilityZone>
            <Name>us-east-1c</Name>
            <ProvisionedIopsCapable>false</ProvisionedIopsCapable>
          </SubnetAvailabilityZone>
        </Subnet>
      </Subnets>
    </DBSubnetGroup>
  </CreateDBSubnetGroupResult>
  <ResponseMetadata>
    <RequestId>3a401b3f-bb9e-11d3-f4c6-37db295f7674</RequestId>
  </ResponseMetadata>
</CreateDBSubnetGroupResponse>
`

var DescribeDBSnapshotsExample = `
<DescribeDBSnapshotsResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <DescribeDBSnapshotsResult>
    <DBSnapshots>
      <DBSnapshot>
        <Port>3306</Port>
        <OptionGroupName>default:mysql-5-6</OptionGroupName>
        <Engine>mysql</Engine>
        <Status>available</Status>
        <SnapshotType>manual</SnapshotType>
        <LicenseModel>general-public-license</LicenseModel>
        <EngineVersion>5.6.13</EngineVersion>
        <DBInstanceIdentifier>my-mysqlexampledb</DBInstanceIdentifier>
        <DBSnapshotIdentifier>my-test-restore-snapshot</DBSnapshotIdentifier>
        <SnapshotCreateTime>2014-03-28T19:57:16.707Z</SnapshotCreateTime>
        <AvailabilityZone>us-west-2b</AvailabilityZone>
        <InstanceCreateTime>2014-01-29T22:58:24.231Z</InstanceCreateTime>
        <PercentProgress>100</PercentProgress>
        <AllocatedStorage>5</AllocatedStorage>
        <MasterUsername>awsmyuser</MasterUsername>
      </DBSnapshot>
      <DBSnapshot>
        <Port>3306</Port>
        <OptionGroupName>default:mysql-5-6</OptionGroupName>
        <Engine>mysql</Engine>
        <Status>available</Status>
        <SnapshotType>automated</SnapshotType>
        <LicenseModel>general-public-license</LicenseModel>
        <EngineVersion>5.6.13</EngineVersion>
        <DBInstanceIdentifier>my-mysqlexampledb</DBInstanceIdentifier>
        <DBSnapshotIdentifier>rds:my-mysqlexampledb-2014-04-19-10-08</DBSnapshotIdentifier>
        <SnapshotCreateTime>2014-04-19T10:09:09.790Z</SnapshotCreateTime>
        <AvailabilityZone>us-west-2b</AvailabilityZone>
        <InstanceCreateTime>2014-01-29T22:58:24.231Z</InstanceCreateTime>
        <PercentProgress>100</PercentProgress>
        <AllocatedStorage>5</AllocatedStorage>
        <MasterUsername>awsmyuser</MasterUsername>
      </DBSnapshot>
      <DBSnapshot>
        <Port>3306</Port>
        <OptionGroupName>default:mysql-5-6</OptionGroupName>
        <Engine>mysql</Engine>
        <Status>available</Status>
        <SnapshotType>automated</SnapshotType>
        <LicenseModel>general-public-license</LicenseModel>
        <EngineVersion>5.6.13</EngineVersion>
        <DBInstanceIdentifier>my-mysqlexampledb</DBInstanceIdentifier>
        <DBSnapshotIdentifier>rds:my-mysqlexampledb-2014-04-20-10-09</DBSnapshotIdentifier>
        <SnapshotCreateTime>2014-04-20T10:09:15.446Z</SnapshotCreateTime>
        <AvailabilityZone>us-west-2b</AvailabilityZone>
        <InstanceCreateTime>2014-01-29T22:58:24.231Z</InstanceCreateTime>
        <PercentProgress>100</PercentProgress>
        <AllocatedStorage>5</AllocatedStorage>
        <MasterUsername>awsmyuser</MasterUsername>
      </DBSnapshot>
    </DBSnapshots>
  </DescribeDBSnapshotsResult>
  <ResponseMetadata>
    <RequestId>b7769930-b98c-11d3-f272-7cd6cce12cc5</RequestId>
  </ResponseMetadata>
</DescribeDBSnapshotsResponse>
`

var RestoreDBInstanceFromDBSnapshotExample = `
<RestoreDBInstanceFromDBSnapshotResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <RestoreDBInstanceFromDBSnapshotResult>
    <DBInstance>
      <BackupRetentionPeriod>7</BackupRetentionPeriod>
      <MultiAZ>false</MultiAZ>
      <DBInstanceStatus>creating</DBInstanceStatus>
      <VpcSecurityGroups/>
      <DBInstanceIdentifier>mysqldb-restored</DBInstanceIdentifier>
      <PreferredBackupWindow>08:14-08:44</PreferredBackupWindow>
      <PreferredMaintenanceWindow>fri:04:50-fri:05:20</PreferredMaintenanceWindow>
      <ReadReplicaDBInstanceIdentifiers/>
      <Engine>mysql</Engine>
      <PendingModifiedValues/>
      <LicenseModel>general-public-license</LicenseModel>
      <EngineVersion>5.6.13</EngineVersion>
      <DBParameterGroups>
        <DBParameterGroup>
          <ParameterApplyStatus>in-sync</ParameterApplyStatus>
          <DBParameterGroupName>default.mysql5.6</DBParameterGroupName>
        </DBParameterGroup>
      </DBParameterGroups>
      <OptionGroupMemberships>
        <OptionGroupMembership>
          <OptionGroupName>default:mysql-5-6</OptionGroupName>
          <Status>pending-apply</Status>
        </OptionGroupMembership>
      </OptionGroupMemberships>
      <PubliclyAccessible>true</PubliclyAccessible>
      <DBSecurityGroups>
        <DBSecurityGroup>
          <Status>active</Status>
          <DBSecurityGroupName>default</DBSecurityGroupName>
        </DBSecurityGroup>
      </DBSecurityGroups>
      <DBName>mysqldb</DBName>
      <AutoMinorVersionUpgrade>true</AutoMinorVersionUpgrade>
      <AllocatedStorage>100</AllocatedStorage>
      <MasterUsername>myawsuser</MasterUsername>
      <DBInstanceClass>db.m1.medium</DBInstanceClass>
    </DBInstance>
  </RestoreDBInstanceFromDBSnapshotResult>
  <ResponseMetadata>
    <RequestId>863fd73e-be2b-11d3-855b-576787000e19</RequestId>
  </ResponseMetadata>
</RestoreDBInstanceFromDBSnapshotResponse>
`

var CreateDBParameterGroupExample = `
<CreateDBParameterGroupResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <CreateDBParameterGroupResult>
    <DBParameterGroup>
      <DBParameterGroupFamily>mysql5.1</DBParameterGroupFamily>
      <Description>My new DB Parameter Group</Description>
      <DBParameterGroupName>mydbparamgroup3</DBParameterGroupName>
    </DBParameterGroup>
  </CreateDBParameterGroupResult>
  <ResponseMetadata>
    <RequestId>7805c127-af22-11c3-96ac-6999cc5f7e72</RequestId>
  </ResponseMetadata>
</CreateDBParameterGroupResponse>
`

var DescribeDBParameterGroupsExample = `
<DescribeDBParameterGroupsResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <DescribeDBParameterGroupsResult>
    <DBParameterGroups>
      <DBParameterGroup>
        <DBParameterGroupFamily>mysql5.6</DBParameterGroupFamily>
        <Description>My new DB Parameter Group</Description>
        <DBParameterGroupName>mydbparamgroup3</DBParameterGroupName>
      </DBParameterGroup>
    </DBParameterGroups>
  </DescribeDBParameterGroupsResult>
  <ResponseMetadata>
    <RequestId>b75d527a-b98c-11d3-f272-7cd6cce12cc5</RequestId>
  </ResponseMetadata>
</DescribeDBParameterGroupsResponse>
`


var DeleteDBParameterGroupExample = `
<DeleteDBParameterGroupResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <ResponseMetadata>
    <RequestId>cad6c267-ba25-11d3-fe11-33d33a9bb7e3</RequestId>
  </ResponseMetadata>
</DeleteDBParameterGroupResponse>
`

var ModifyDBParameterGroupExample = `
<ModifyDBParameterGroupResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <ModifyDBParameterGroupResult>
    <DBParameterGroupName>mydbparamgroup3</DBParameterGroupName>
  </ModifyDBParameterGroupResult>
  <ResponseMetadata>
    <RequestId>12d7435e-bba0-11d3-fe11-33d33a9bb7e3</RequestId>
  </ResponseMetadata>
</ModifyDBParameterGroupResponse>
`

var DescribeDBParametersExample = `
<DescribeDBParametersResponse xmlns="http://rds.amazonaws.com/doc/2014-09-01/">
  <DescribeDBParametersResult>
    <Marker>bGlzdGVuZXJfbmV0d29ya3M=</Marker>
    <Parameters>
      <Parameter>
        <ParameterValue>utf8</ParameterValue>
        <ParameterName>character_set_server</ParameterName>
      </Parameter>
      <Parameter>
        <ParameterValue>utf8</ParameterValue>
        <ParameterName>character_set_client</ParameterName>
      </Parameter>
      <Parameter>
        <ParameterValue>utf8</ParameterValue>
        <ParameterName>character_set_results</ParameterName>
      </Parameter>
      <Parameter>
        <ParameterValue>utf8_unicode_ci</ParameterValue>
        <ParameterName>collation_server</ParameterName>
      </Parameter>
      <Parameter>
        <ParameterValue>utf8_unicode_ci</ParameterValue>
        <ParameterName>collation_connection</ParameterName>
      </Parameter>
    </Parameters>
  </DescribeDBParametersResult>
  <ResponseMetadata>
    <RequestId>8c40488f-b9ff-11d3-a15e-7ac49293f4fa</RequestId>
  </ResponseMetadata>
</DescribeDBParametersResponse>
`
