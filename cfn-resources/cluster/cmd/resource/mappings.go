// Copyright 2023 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20231001001/admin"
)

func mapClusterToModel(model *Model, cluster *admin.AdvancedClusterDescription) {
	model.Id = cluster.Id
	model.ProjectId = cluster.GroupId
	model.Name = cluster.Name
	model.BackupEnabled = cluster.BackupEnabled
	model.BiConnector = flattenBiConnectorConfig(cluster.BiConnector)
	model.ConnectionStrings = flattenConnectionStrings(cluster.ConnectionStrings)
	model.ClusterType = cluster.ClusterType
	model.CreatedDate = util.TimePtrToStringPtr(cluster.CreateDate)
	model.DiskSizeGB = cluster.DiskSizeGB
	model.EncryptionAtRestProvider = cluster.EncryptionAtRestProvider
	model.Labels = flattenLabels(cluster.Labels)
	model.MongoDBMajorVersion = cluster.MongoDBMajorVersion
	model.MongoDBVersion = cluster.MongoDBVersion
	model.Paused = cluster.Paused
	model.PitEnabled = cluster.PitEnabled
	model.RootCertType = cluster.RootCertType
	model.ReplicationSpecs = flattenReplicationSpecs(cluster.ReplicationSpecs)
	model.StateName = cluster.StateName
	model.VersionReleaseSystem = cluster.VersionReleaseSystem
}

func containsLabelOrKey(list []Labels, item Labels) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, item) || *v.Key == *item.Key {
			return true
		}
	}

	return false
}

func expandBiConnector(biConnector *BiConnector) *admin.BiConnector {
	if biConnector == nil {
		return nil
	}
	return &admin.BiConnector{
		Enabled:        biConnector.Enabled,
		ReadPreference: biConnector.ReadPreference,
	}
}

func expandReplicationSpecs(replicationSpecs []AdvancedReplicationSpec) []admin.ReplicationSpec {
	var rSpecs []admin.ReplicationSpec

	for i := range replicationSpecs {
		var numShards int

		rSpec := admin.ReplicationSpec{
			NumShards:     &numShards,
			RegionConfigs: expandRegionsConfig(replicationSpecs[i].AdvancedRegionConfigs),
		}

		if util.IsStringPresent(replicationSpecs[i].ID) {
			rSpec.Id = admin.PtrString(cast.ToString(replicationSpecs[i].ID))
		}

		if replicationSpecs[i].NumShards != nil {
			rSpec.NumShards = replicationSpecs[i].NumShards
		}
		if replicationSpecs[i].ZoneName != nil {
			rSpec.ZoneName = admin.PtrString(cast.ToString(replicationSpecs[i].ZoneName))
		}
		rSpecs = append(rSpecs, rSpec)
	}

	fmt.Printf("specs: len %d %+v", len(replicationSpecs), rSpecs)
	return rSpecs
}

func expandAutoScaling(scaling *AdvancedAutoScaling) *admin.AdvancedAutoScalingSettings {
	advAutoScaling := &admin.AdvancedAutoScalingSettings{}
	if scaling == nil {
		return nil
	}
	if scaling.Compute != nil {
		advAutoScaling.Compute = &admin.AdvancedComputeAutoScaling{
			Enabled:          scaling.Compute.Enabled,
			ScaleDownEnabled: scaling.Compute.ScaleDownEnabled,
		}

		if util.IsStringPresent(scaling.Compute.MinInstanceSize) {
			advAutoScaling.Compute.MinInstanceSize = scaling.Compute.MinInstanceSize
		}

		if util.IsStringPresent(scaling.Compute.MaxInstanceSize) {
			advAutoScaling.Compute.MaxInstanceSize = scaling.Compute.MaxInstanceSize
		}
	}

	if scaling.DiskGB != nil {
		advAutoScaling.DiskGB = &admin.DiskGBAutoScaling{Enabled: scaling.DiskGB.Enabled}
	}

	return advAutoScaling
}

func expandRegionsConfig(regionConfigs []AdvancedRegionConfig) []admin.CloudRegionConfig {
	var regionsConfigs []admin.CloudRegionConfig
	for _, regionCfg := range regionConfigs {
		regionsConfigs = append(regionsConfigs, expandRegionConfig(regionCfg))
	}
	return regionsConfigs
}

func expandRegionConfig(regionCfg AdvancedRegionConfig) admin.CloudRegionConfig {
	var region string
	if regionCfg.RegionName != nil {
		region = *regionCfg.RegionName
	}

	providerName := constants.AWS
	if regionCfg.ProviderName != nil {
		providerName = *regionCfg.ProviderName
	}

	advRegionConfig := admin.CloudRegionConfig{
		ProviderName: &providerName,
		RegionName:   &region,
		Priority:     regionCfg.Priority,
	}

	if regionCfg.AutoScaling != nil {
		advRegionConfig.AutoScaling = expandAutoScaling(regionCfg.AutoScaling)
	}
	if regionCfg.AnalyticsAutoScaling != nil {
		advRegionConfig.AnalyticsAutoScaling = expandAutoScaling(regionCfg.AnalyticsAutoScaling)
	}
	if regionCfg.AnalyticsSpecs != nil {
		advRegionConfig.AnalyticsSpecs = expandRegionConfigSpec(regionCfg.AnalyticsSpecs)
	}
	if regionCfg.ElectableSpecs != nil {
		advRegionConfig.ElectableSpecs = newHardwareSpec(regionCfg.ElectableSpecs)
	}
	if regionCfg.ReadOnlySpecs != nil {
		advRegionConfig.ReadOnlySpecs = expandRegionConfigSpec(regionCfg.ReadOnlySpecs)
	}
	if regionCfg.BackingProviderName != nil {
		advRegionConfig.BackingProviderName = regionCfg.BackingProviderName
	}
	return advRegionConfig
}

func newHardwareSpec(spec *Specs) *admin.HardwareSpec {
	if spec == nil {
		return nil
	}

	var ebsVolumeType string
	var instanceSize string
	if spec.EbsVolumeType != nil {
		ebsVolumeType = *spec.EbsVolumeType
	}
	if spec.InstanceSize != nil {
		instanceSize = *spec.InstanceSize
	}
	var val int64
	if spec.DiskIOPS != nil {
		v, err := strconv.ParseInt(*spec.DiskIOPS, 10, 64)
		if err == nil {
			val = v
		}
	}

	return &admin.HardwareSpec{
		DiskIOPS:      util.Int64PtrToIntPtr(&val),
		EbsVolumeType: &ebsVolumeType,
		InstanceSize:  &instanceSize,
		NodeCount:     spec.NodeCount,
	}
}

func expandRegionConfigSpec(spec *Specs) *admin.DedicatedHardwareSpec {
	if spec == nil {
		return nil
	}

	var ebsVolumeType string
	var instanceSize string
	if spec.EbsVolumeType != nil {
		ebsVolumeType = *spec.EbsVolumeType
	}
	if spec.InstanceSize != nil {
		instanceSize = *spec.InstanceSize
	}
	var val int64
	if spec.DiskIOPS != nil {
		v, err := strconv.ParseInt(*spec.DiskIOPS, 10, 64)
		if err == nil {
			val = v
		}
	}

	return &admin.DedicatedHardwareSpec{
		DiskIOPS:      util.Int64PtrToIntPtr(&val),
		EbsVolumeType: &ebsVolumeType,
		InstanceSize:  &instanceSize,
		NodeCount:     spec.NodeCount,
	}
}

func expandLabelSlice(labels []Labels) []admin.ComponentLabel {
	res := make([]admin.ComponentLabel, len(labels))

	for i := range labels {
		var key string
		if labels[i].Key != nil {
			key = *labels[i].Key
		}
		var value string
		if labels[i].Key != nil {
			value = *labels[i].Value
		}
		res[i] = admin.ComponentLabel{
			Key:   &key,
			Value: &value,
		}
	}
	return res
}

func flattenAutoScaling(scaling *admin.AdvancedAutoScalingSettings) *AdvancedAutoScaling {
	if scaling == nil {
		return nil
	}
	advAutoScaling := &AdvancedAutoScaling{}

	if scaling.DiskGB != nil {
		advAutoScaling.DiskGB = &DiskGB{Enabled: scaling.DiskGB.Enabled}
	}
	if scaling.Compute != nil {
		compute := &Compute{}
		if scaling.Compute.Enabled != nil {
			compute.Enabled = scaling.Compute.Enabled
		}
		if scaling.Compute.ScaleDownEnabled != nil {
			compute.ScaleDownEnabled = scaling.Compute.ScaleDownEnabled
		}
		if util.IsStringPresent(scaling.Compute.MinInstanceSize) {
			compute.MinInstanceSize = scaling.Compute.MinInstanceSize
		}
		if util.IsStringPresent(scaling.Compute.MaxInstanceSize) {
			compute.MaxInstanceSize = scaling.Compute.MaxInstanceSize
		}

		advAutoScaling.Compute = compute
	}
	return advAutoScaling
}

func flattenReplicationSpecs(replicationSpecs []admin.ReplicationSpec) []AdvancedReplicationSpec {
	var rSpecs []AdvancedReplicationSpec

	for ind := range replicationSpecs {
		rSpec := AdvancedReplicationSpec{
			ID:                    replicationSpecs[ind].Id,
			NumShards:             replicationSpecs[ind].NumShards,
			ZoneName:              replicationSpecs[ind].ZoneName,
			AdvancedRegionConfigs: flattenRegionsConfig(replicationSpecs[ind].RegionConfigs),
		}
		rSpecs = append(rSpecs, rSpec)
	}
	fmt.Printf("specs: len %d %+v", len(replicationSpecs), rSpecs)
	return rSpecs
}

func flattenRegionsConfig(regionConfigs []admin.CloudRegionConfig) []AdvancedRegionConfig {
	var regionsConfigs []AdvancedRegionConfig
	for i := range regionConfigs {
		regionsConfigs = append(regionsConfigs, flattenRegionConfig(&regionConfigs[i]))
	}
	return regionsConfigs
}

func flattenRegionConfig(regionCfg *admin.CloudRegionConfig) AdvancedRegionConfig {
	advRegConfig := AdvancedRegionConfig{
		AutoScaling:          flattenAutoScaling(regionCfg.AutoScaling),
		AnalyticsAutoScaling: flattenAutoScaling(regionCfg.AnalyticsAutoScaling),
		RegionName:           regionCfg.RegionName,
		Priority:             regionCfg.Priority,
	}
	if regionCfg.AnalyticsSpecs != nil {
		advRegConfig.AnalyticsSpecs = flattenRegionConfigSpec(regionCfg.AnalyticsSpecs)
	}
	if regionCfg.ElectableSpecs != nil {
		advRegConfig.ElectableSpecs = flattenElectableSpecs(regionCfg.ElectableSpecs)
	}

	if regionCfg.ReadOnlySpecs != nil {
		advRegConfig.ReadOnlySpecs = flattenRegionConfigSpec(regionCfg.ReadOnlySpecs)
	}

	return advRegConfig
}

func flattenElectableSpecs(spec *admin.HardwareSpec) *Specs {
	if spec == nil {
		return nil
	}
	var diskIops string
	if spec.DiskIOPS != nil {
		diskIops = strconv.Itoa(*spec.DiskIOPS)
	}

	return &Specs{
		DiskIOPS:      &diskIops,
		EbsVolumeType: spec.EbsVolumeType,
		InstanceSize:  spec.InstanceSize,
		NodeCount:     spec.NodeCount,
	}
}

func flattenRegionConfigSpec(spec *admin.DedicatedHardwareSpec) *Specs {
	if spec == nil {
		return nil
	}
	var diskIops string
	if spec.DiskIOPS != nil {
		diskIops = strconv.Itoa(*spec.DiskIOPS)
	}

	return &Specs{
		DiskIOPS:      &diskIops,
		EbsVolumeType: spec.EbsVolumeType,
		InstanceSize:  spec.InstanceSize,
		NodeCount:     spec.NodeCount,
	}
}

func flattenBiConnectorConfig(biConnector *admin.BiConnector) *BiConnector {
	if biConnector == nil {
		return nil
	}

	return &BiConnector{
		ReadPreference: biConnector.ReadPreference,
		Enabled:        biConnector.Enabled,
	}
}

type privateEndpointConnectionStrings struct {
	PrivateEndpoints                  []string
	PrivateEndpointsSrv               []string
	SRVShardOptimizedConnectionString []string
}

func flattenConnectionStrings(clusterConnStrings *admin.ClusterConnectionStrings) (connStrings *ConnectionStrings) {
	if clusterConnStrings != nil {
		privateEndpoints := flattenPrivateEndpoint(clusterConnStrings.PrivateEndpoint)
		connStrings = &ConnectionStrings{
			Standard:                          clusterConnStrings.Standard,
			StandardSrv:                       clusterConnStrings.StandardSrv,
			Private:                           clusterConnStrings.Private,
			PrivateSrv:                        clusterConnStrings.PrivateSrv,
			PrivateEndpoints:                  privateEndpoints.PrivateEndpoints,
			PrivateEndpointsSrv:               privateEndpoints.PrivateEndpointsSrv,
			SRVShardOptimizedConnectionString: privateEndpoints.SRVShardOptimizedConnectionString,
		}
	}
	return
}

func flattenPrivateEndpoint(pes []admin.ClusterDescriptionConnectionStringsPrivateEndpoint) privateEndpointConnectionStrings {
	privateEndpoints := privateEndpointConnectionStrings{
		PrivateEndpoints:                  make([]string, 0),
		PrivateEndpointsSrv:               make([]string, 0),
		SRVShardOptimizedConnectionString: make([]string, 0),
	}

	for _, pe := range pes {
		if util.IsStringPresent(pe.ConnectionString) {
			privateEndpoints.PrivateEndpoints = append(privateEndpoints.PrivateEndpoints, *pe.ConnectionString)
		}

		if util.IsStringPresent(pe.SrvConnectionString) {
			privateEndpoints.PrivateEndpointsSrv = append(privateEndpoints.PrivateEndpointsSrv, *pe.SrvConnectionString)
		}

		if util.IsStringPresent(pe.SrvShardOptimizedConnectionString) {
			privateEndpoints.SRVShardOptimizedConnectionString = append(privateEndpoints.SRVShardOptimizedConnectionString, *pe.SrvShardOptimizedConnectionString)
		}
	}
	return privateEndpoints
}

func flattenProcessArgs(p *admin.ClusterDescriptionProcessArgs) *ProcessArgs {
	return &ProcessArgs{
		DefaultReadConcern:               p.DefaultReadConcern,
		DefaultWriteConcern:              p.DefaultWriteConcern,
		FailIndexKeyTooLong:              p.FailIndexKeyTooLong,
		JavascriptEnabled:                p.JavascriptEnabled,
		MinimumEnabledTLSProtocol:        p.MinimumEnabledTlsProtocol,
		NoTableScan:                      p.NoTableScan,
		OplogSizeMB:                      p.OplogSizeMB,
		SampleSizeBIConnector:            p.SampleSizeBIConnector,
		SampleRefreshIntervalBIConnector: p.SampleRefreshIntervalBIConnector,
		OplogMinRetentionHours:           p.OplogMinRetentionHours,
		TransactionLifetimeLimitSeconds:  util.Int64PtrToIntPtr(p.TransactionLifetimeLimitSeconds),
	}
}

func flattenLabels(clusterLabels []admin.ComponentLabel) []Labels {
	labels := make([]Labels, len(clusterLabels))
	for i := range clusterLabels {
		labels[i] = Labels{
			Key:   clusterLabels[i].Key,
			Value: clusterLabels[i].Value,
		}
	}
	return labels
}

func expandAdvancedSettings(processArgs ProcessArgs) *admin.ClusterDescriptionProcessArgs {
	var args admin.ClusterDescriptionProcessArgs

	if processArgs.DefaultReadConcern != nil {
		args.DefaultReadConcern = processArgs.DefaultReadConcern
	}
	args.FailIndexKeyTooLong = processArgs.FailIndexKeyTooLong
	if processArgs.DefaultWriteConcern != nil {
		args.DefaultWriteConcern = processArgs.DefaultWriteConcern
	}
	args.JavascriptEnabled = processArgs.JavascriptEnabled
	if processArgs.MinimumEnabledTLSProtocol != nil {
		args.MinimumEnabledTlsProtocol = processArgs.MinimumEnabledTLSProtocol
	}
	args.NoTableScan = processArgs.NoTableScan

	if processArgs.OplogSizeMB != nil {
		args.OplogSizeMB = processArgs.OplogSizeMB
	}
	if processArgs.SampleSizeBIConnector != nil {
		args.SampleSizeBIConnector = processArgs.SampleSizeBIConnector
	}
	if processArgs.SampleRefreshIntervalBIConnector != nil {
		args.SampleRefreshIntervalBIConnector = processArgs.SampleRefreshIntervalBIConnector
	}

	if processArgs.OplogMinRetentionHours != nil {
		args.OplogMinRetentionHours = processArgs.OplogMinRetentionHours
	}

	if processArgs.TransactionLifetimeLimitSeconds != nil {
		args.TransactionLifetimeLimitSeconds = cast64(processArgs.TransactionLifetimeLimitSeconds)
	}

	return &args
}

func flattenTags(clusterTags []admin.ResourceTag) (tags []Tag) {
	for ind := range clusterTags {
		tags = append(tags, Tag{
			Key:   clusterTags[ind].Key,
			Value: clusterTags[ind].Value,
		})
	}
	return
}

func expandTags(tags []Tag) (clusterTags []admin.ResourceTag) {
	for ind := range tags {
		clusterTags = append(clusterTags, admin.ResourceTag{
			Key:   tags[ind].Key,
			Value: tags[ind].Value,
		})
	}
	return
}

func setClusterData(currentModel *Model, cluster *admin.AdvancedClusterDescription) {
	if cluster == nil {
		return
	}

	currentModel.ProjectId = cluster.GroupId
	currentModel.Name = cluster.Name
	currentModel.Id = cluster.Id

	if currentModel.BackupEnabled != nil {
		currentModel.BackupEnabled = cluster.BackupEnabled
	}
	if currentModel.BiConnector != nil {
		currentModel.BiConnector = flattenBiConnectorConfig(cluster.BiConnector)
	}
	// Readonly
	currentModel.ConnectionStrings = flattenConnectionStrings(cluster.ConnectionStrings)
	if currentModel.ClusterType != nil {
		currentModel.ClusterType = cluster.ClusterType
	}
	// Readonly
	currentModel.CreatedDate = util.TimePtrToStringPtr(cluster.CreateDate)
	if currentModel.DiskSizeGB != nil {
		currentModel.DiskSizeGB = cluster.DiskSizeGB
	}
	if currentModel.EncryptionAtRestProvider != nil {
		currentModel.EncryptionAtRestProvider = cluster.EncryptionAtRestProvider
	}
	if currentModel.Labels != nil {
		currentModel.Labels = flattenLabels(cluster.Labels)
	}
	if currentModel.MongoDBMajorVersion != nil {
		currentModel.MongoDBMajorVersion = cluster.MongoDBMajorVersion
	}
	// Readonly
	currentModel.MongoDBVersion = cluster.MongoDBVersion

	if currentModel.Paused != nil {
		currentModel.Paused = cluster.Paused
	}
	if currentModel.PitEnabled != nil {
		currentModel.PitEnabled = cluster.PitEnabled
	}
	if currentModel.RootCertType != nil {
		currentModel.RootCertType = cluster.RootCertType
	}
	if currentModel.ReplicationSpecs != nil {
		currentModel.ReplicationSpecs = flattenReplicationSpecs(cluster.ReplicationSpecs)
	}
	// Readonly
	currentModel.StateName = cluster.StateName
	if currentModel.VersionReleaseSystem != nil {
		currentModel.VersionReleaseSystem = cluster.VersionReleaseSystem
	}

	currentModel.TerminationProtectionEnabled = cluster.TerminationProtectionEnabled
	currentModel.Tags = flattenTags(cluster.Tags)
}

func setClusterRequest(currentModel *Model) (*admin.AdvancedClusterDescription, handler.ProgressEvent, error) {
	// Atlas client
	clusterRequest := &admin.AdvancedClusterDescription{
		Name:             currentModel.Name,
		ReplicationSpecs: expandReplicationSpecs(currentModel.ReplicationSpecs),
	}

	if currentModel.EncryptionAtRestProvider != nil {
		clusterRequest.EncryptionAtRestProvider = currentModel.EncryptionAtRestProvider
	}

	if currentModel.ClusterType != nil {
		clusterRequest.ClusterType = currentModel.ClusterType
	}

	if currentModel.BackupEnabled != nil {
		clusterRequest.BackupEnabled = currentModel.BackupEnabled
	}

	if currentModel.BiConnector != nil {
		clusterRequest.BiConnector = expandBiConnector(currentModel.BiConnector)
	}

	if currentModel.DiskSizeGB != nil {
		clusterRequest.DiskSizeGB = currentModel.DiskSizeGB
	}

	if len(currentModel.Labels) > 0 {
		clusterRequest.Labels = expandLabelSlice(currentModel.Labels)
	}

	if currentModel.MongoDBMajorVersion != nil {
		clusterRequest.MongoDBMajorVersion = admin.PtrString(formatMongoDBMajorVersion(*currentModel.MongoDBMajorVersion))
	}

	if currentModel.PitEnabled != nil {
		clusterRequest.PitEnabled = currentModel.PitEnabled
	}

	if currentModel.VersionReleaseSystem != nil {
		clusterRequest.VersionReleaseSystem = currentModel.VersionReleaseSystem
	}

	if currentModel.RootCertType != nil {
		clusterRequest.RootCertType = currentModel.RootCertType
	}
	clusterRequest.Tags = expandTags(currentModel.Tags)

	clusterRequest.TerminationProtectionEnabled = currentModel.TerminationProtectionEnabled
	return clusterRequest, handler.ProgressEvent{}, nil
}
