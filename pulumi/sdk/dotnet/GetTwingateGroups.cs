// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;
using Pulumi;

namespace TwingateLabs.Twingate
{
    public static class GetTwingateGroups
    {
        public static Task<GetTwingateGroupsResult> InvokeAsync(GetTwingateGroupsArgs? args = null, InvokeOptions? options = null)
            => global::Pulumi.Deployment.Instance.InvokeAsync<GetTwingateGroupsResult>("twingate:index/getTwingateGroups:getTwingateGroups", args ?? new GetTwingateGroupsArgs(), options.WithDefaults());

        public static Output<GetTwingateGroupsResult> Invoke(GetTwingateGroupsInvokeArgs? args = null, InvokeOptions? options = null)
            => global::Pulumi.Deployment.Instance.Invoke<GetTwingateGroupsResult>("twingate:index/getTwingateGroups:getTwingateGroups", args ?? new GetTwingateGroupsInvokeArgs(), options.WithDefaults());
    }


    public sealed class GetTwingateGroupsArgs : global::Pulumi.InvokeArgs
    {
        [Input("groups")]
        private List<Inputs.GetTwingateGroupsGroupArgs>? _groups;
        public List<Inputs.GetTwingateGroupsGroupArgs> Groups
        {
            get => _groups ?? (_groups = new List<Inputs.GetTwingateGroupsGroupArgs>());
            set => _groups = value;
        }

        [Input("isActive")]
        public bool? IsActive { get; set; }

        [Input("name")]
        public string? Name { get; set; }

        [Input("type")]
        public string? Type { get; set; }

        public GetTwingateGroupsArgs()
        {
        }
        public static new GetTwingateGroupsArgs Empty => new GetTwingateGroupsArgs();
    }

    public sealed class GetTwingateGroupsInvokeArgs : global::Pulumi.InvokeArgs
    {
        [Input("groups")]
        private InputList<Inputs.GetTwingateGroupsGroupInputArgs>? _groups;
        public InputList<Inputs.GetTwingateGroupsGroupInputArgs> Groups
        {
            get => _groups ?? (_groups = new InputList<Inputs.GetTwingateGroupsGroupInputArgs>());
            set => _groups = value;
        }

        [Input("isActive")]
        public Input<bool>? IsActive { get; set; }

        [Input("name")]
        public Input<string>? Name { get; set; }

        [Input("type")]
        public Input<string>? Type { get; set; }

        public GetTwingateGroupsInvokeArgs()
        {
        }
        public static new GetTwingateGroupsInvokeArgs Empty => new GetTwingateGroupsInvokeArgs();
    }


    [OutputType]
    public sealed class GetTwingateGroupsResult
    {
        public readonly ImmutableArray<Outputs.GetTwingateGroupsGroupResult> Groups;
        /// <summary>
        /// The provider-assigned unique ID for this managed resource.
        /// </summary>
        public readonly string Id;
        public readonly bool? IsActive;
        public readonly string? Name;
        public readonly string? Type;

        [OutputConstructor]
        private GetTwingateGroupsResult(
            ImmutableArray<Outputs.GetTwingateGroupsGroupResult> groups,

            string id,

            bool? isActive,

            string? name,

            string? type)
        {
            Groups = groups;
            Id = id;
            IsActive = isActive;
            Name = name;
            Type = type;
        }
    }
}