// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;
using Pulumi;

namespace TwingateLabs.Twingate.Outputs
{

    [OutputType]
    public sealed class GetTwingateResourceProtocolResult
    {
        /// <summary>
        /// Whether to allow ICMP (ping) traffic
        /// </summary>
        public readonly bool AllowIcmp;
        public readonly ImmutableArray<Outputs.GetTwingateResourceProtocolTcpResult> Tcps;
        public readonly ImmutableArray<Outputs.GetTwingateResourceProtocolUdpResult> Udps;

        [OutputConstructor]
        private GetTwingateResourceProtocolResult(
            bool allowIcmp,

            ImmutableArray<Outputs.GetTwingateResourceProtocolTcpResult> tcps,

            ImmutableArray<Outputs.GetTwingateResourceProtocolUdpResult> udps)
        {
            AllowIcmp = allowIcmp;
            Tcps = tcps;
            Udps = udps;
        }
    }
}