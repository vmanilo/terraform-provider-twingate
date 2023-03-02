// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;
using Pulumi;

namespace TwingateLabs.Twingate.Inputs
{

    public sealed class GetTwingateConnectorsConnectorArgs : global::Pulumi.InvokeArgs
    {
        /// <summary>
        /// The ID of the Connector
        /// </summary>
        [Input("id", required: true)]
        public string Id { get; set; } = null!;

        /// <summary>
        /// The Name of the Connector
        /// </summary>
        [Input("name", required: true)]
        public string Name { get; set; } = null!;

        /// <summary>
        /// The ID of the Remote Network attached to the Connector
        /// </summary>
        [Input("remoteNetworkId", required: true)]
        public string RemoteNetworkId { get; set; } = null!;

        public GetTwingateConnectorsConnectorArgs()
        {
        }
        public static new GetTwingateConnectorsConnectorArgs Empty => new GetTwingateConnectorsConnectorArgs();
    }
}