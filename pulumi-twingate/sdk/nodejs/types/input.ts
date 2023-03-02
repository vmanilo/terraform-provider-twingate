// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../types/input";
import * as outputs from "../types/output";

export interface GetTwingateConnectorsConnector {
    /**
     * The ID of the Connector
     */
    id?: string;
    /**
     * The Name of the Connector
     */
    name?: string;
    /**
     * The ID of the Remote Network attached to the Connector
     */
    remoteNetworkId?: string;
}

export interface GetTwingateConnectorsConnectorArgs {
    /**
     * The ID of the Connector
     */
    id?: pulumi.Input<string>;
    /**
     * The Name of the Connector
     */
    name?: pulumi.Input<string>;
    /**
     * The ID of the Remote Network attached to the Connector
     */
    remoteNetworkId?: pulumi.Input<string>;
}

export interface GetTwingateGroupsGroup {
    /**
     * The ID of the Group
     */
    id?: string;
    /**
     * Indicates if the Group is active
     */
    isActive?: boolean;
    /**
     * The name of the Group
     */
    name?: string;
    /**
     * The type of the Group
     */
    type?: string;
}

export interface GetTwingateGroupsGroupArgs {
    /**
     * The ID of the Group
     */
    id?: pulumi.Input<string>;
    /**
     * Indicates if the Group is active
     */
    isActive?: pulumi.Input<boolean>;
    /**
     * The name of the Group
     */
    name?: pulumi.Input<string>;
    /**
     * The type of the Group
     */
    type?: pulumi.Input<string>;
}

export interface GetTwingateRemoteNetworksRemoteNetwork {
    /**
     * The ID of the Remote Network
     */
    id?: string;
    /**
     * The location of the Remote Network. Must be one of the following: AWS, AZURE, GOOGLE*CLOUD, ON*PREMISE, OTHER.
     */
    location?: string;
    /**
     * The name of the Remote Network
     */
    name?: string;
}

export interface GetTwingateRemoteNetworksRemoteNetworkArgs {
    /**
     * The ID of the Remote Network
     */
    id?: pulumi.Input<string>;
    /**
     * The location of the Remote Network. Must be one of the following: AWS, AZURE, GOOGLE*CLOUD, ON*PREMISE, OTHER.
     */
    location?: pulumi.Input<string>;
    /**
     * The name of the Remote Network
     */
    name?: pulumi.Input<string>;
}

export interface GetTwingateResourceProtocol {
    /**
     * Whether to allow ICMP (ping) traffic
     */
    allowIcmp?: boolean;
    tcps?: inputs.GetTwingateResourceProtocolTcp[];
    udps?: inputs.GetTwingateResourceProtocolUdp[];
}

export interface GetTwingateResourceProtocolArgs {
    /**
     * Whether to allow ICMP (ping) traffic
     */
    allowIcmp?: pulumi.Input<boolean>;
    tcps?: pulumi.Input<pulumi.Input<inputs.GetTwingateResourceProtocolTcpArgs>[]>;
    udps?: pulumi.Input<pulumi.Input<inputs.GetTwingateResourceProtocolUdpArgs>[]>;
}

export interface GetTwingateResourceProtocolTcp {
    policy?: string;
    ports?: string[];
}

export interface GetTwingateResourceProtocolTcpArgs {
    policy?: pulumi.Input<string>;
    ports?: pulumi.Input<pulumi.Input<string>[]>;
}

export interface GetTwingateResourceProtocolUdp {
    policy?: string;
    ports?: string[];
}

export interface GetTwingateResourceProtocolUdpArgs {
    policy?: pulumi.Input<string>;
    ports?: pulumi.Input<pulumi.Input<string>[]>;
}

export interface GetTwingateResourcesResource {
    /**
     * The Resource's IP/CIDR or FQDN/DNS zone
     */
    address?: string;
    /**
     * The id of the Resource
     */
    id?: string;
    /**
     * The name of the Resource
     */
    name?: string;
    /**
     * Restrict access to certain protocols and ports. By default or when this argument is not defined, there is no restriction, and all protocols and ports are allowed.
     */
    protocols?: inputs.GetTwingateResourcesResourceProtocol[];
    /**
     * Remote Network ID where the Resource lives
     */
    remoteNetworkId?: string;
}

export interface GetTwingateResourcesResourceArgs {
    /**
     * The Resource's IP/CIDR or FQDN/DNS zone
     */
    address?: pulumi.Input<string>;
    /**
     * The id of the Resource
     */
    id?: pulumi.Input<string>;
    /**
     * The name of the Resource
     */
    name?: pulumi.Input<string>;
    /**
     * Restrict access to certain protocols and ports. By default or when this argument is not defined, there is no restriction, and all protocols and ports are allowed.
     */
    protocols?: pulumi.Input<pulumi.Input<inputs.GetTwingateResourcesResourceProtocolArgs>[]>;
    /**
     * Remote Network ID where the Resource lives
     */
    remoteNetworkId?: pulumi.Input<string>;
}

export interface GetTwingateResourcesResourceProtocol {
    allowIcmp?: boolean;
    tcps?: inputs.GetTwingateResourcesResourceProtocolTcp[];
    udps?: inputs.GetTwingateResourcesResourceProtocolUdp[];
}

export interface GetTwingateResourcesResourceProtocolArgs {
    allowIcmp?: pulumi.Input<boolean>;
    tcps?: pulumi.Input<pulumi.Input<inputs.GetTwingateResourcesResourceProtocolTcpArgs>[]>;
    udps?: pulumi.Input<pulumi.Input<inputs.GetTwingateResourcesResourceProtocolUdpArgs>[]>;
}

export interface GetTwingateResourcesResourceProtocolTcp {
    policy?: string;
    ports?: string[];
}

export interface GetTwingateResourcesResourceProtocolTcpArgs {
    policy?: pulumi.Input<string>;
    ports?: pulumi.Input<pulumi.Input<string>[]>;
}

export interface GetTwingateResourcesResourceProtocolUdp {
    policy?: string;
    ports?: string[];
}

export interface GetTwingateResourcesResourceProtocolUdpArgs {
    policy?: pulumi.Input<string>;
    ports?: pulumi.Input<pulumi.Input<string>[]>;
}

export interface GetTwingateSecurityPoliciesSecurityPolicy {
    /**
     * Return a matching Security Policy by its ID. The ID for the Security Policy must be obtained from the Admin API.
     */
    id?: string;
    /**
     * Return a Security Policy that exactly matches this name.
     */
    name?: string;
}

export interface GetTwingateSecurityPoliciesSecurityPolicyArgs {
    /**
     * Return a matching Security Policy by its ID. The ID for the Security Policy must be obtained from the Admin API.
     */
    id?: pulumi.Input<string>;
    /**
     * Return a Security Policy that exactly matches this name.
     */
    name?: pulumi.Input<string>;
}

export interface GetTwingateServiceAccountsServiceAccount {
    /**
     * ID of the Service Account resource
     */
    id?: string;
    /**
     * List of twingate*service*account_key IDs that are assigned to the Service Account.
     */
    keyIds?: string[];
    /**
     * Name of the Service Account
     */
    name?: string;
    /**
     * List of twingate.TwingateResource IDs that the Service Account is assigned to.
     */
    resourceIds?: string[];
}

export interface GetTwingateServiceAccountsServiceAccountArgs {
    /**
     * ID of the Service Account resource
     */
    id?: pulumi.Input<string>;
    /**
     * List of twingate*service*account_key IDs that are assigned to the Service Account.
     */
    keyIds?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * Name of the Service Account
     */
    name?: pulumi.Input<string>;
    /**
     * List of twingate.TwingateResource IDs that the Service Account is assigned to.
     */
    resourceIds?: pulumi.Input<pulumi.Input<string>[]>;
}

export interface GetTwingateUsersUser {
    /**
     * The email address of the User
     */
    email?: string;
    /**
     * The first name of the User
     */
    firstName?: string;
    /**
     * The ID of the User
     */
    id?: string;
    /**
     * Indicates whether the User is an admin
     */
    isAdmin?: boolean;
    /**
     * The last name of the User
     */
    lastName?: string;
    /**
     * Indicates the User's role. Either ADMIN, DEVOPS, SUPPORT, or MEMBER.
     */
    role?: string;
}

export interface GetTwingateUsersUserArgs {
    /**
     * The email address of the User
     */
    email?: pulumi.Input<string>;
    /**
     * The first name of the User
     */
    firstName?: pulumi.Input<string>;
    /**
     * The ID of the User
     */
    id?: pulumi.Input<string>;
    /**
     * Indicates whether the User is an admin
     */
    isAdmin?: pulumi.Input<boolean>;
    /**
     * The last name of the User
     */
    lastName?: pulumi.Input<string>;
    /**
     * Indicates the User's role. Either ADMIN, DEVOPS, SUPPORT, or MEMBER.
     */
    role?: pulumi.Input<string>;
}

export interface TwingateResourceAccess {
    /**
     * List of Group IDs that will have permission to access the Resource.
     */
    groupIds?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * List of Service Account IDs that will have permission to access the Resource.
     */
    serviceAccountIds?: pulumi.Input<pulumi.Input<string>[]>;
}

export interface TwingateResourceProtocols {
    /**
     * Whether to allow ICMP (ping) traffic
     */
    allowIcmp?: pulumi.Input<boolean>;
    tcp: pulumi.Input<inputs.TwingateResourceProtocolsTcp>;
    udp: pulumi.Input<inputs.TwingateResourceProtocolsUdp>;
}

export interface TwingateResourceProtocolsTcp {
    policy: pulumi.Input<string>;
    ports?: pulumi.Input<pulumi.Input<string>[]>;
}

export interface TwingateResourceProtocolsUdp {
    policy: pulumi.Input<string>;
    ports?: pulumi.Input<pulumi.Input<string>[]>;
}